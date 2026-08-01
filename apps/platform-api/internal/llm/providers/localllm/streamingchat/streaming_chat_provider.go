package streamingchat

import (
	"bufio"
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/go-packages/json"
	streamingchatPkg "github.com/blocknextai/platform-api/internal/llm/streamingchat"
)

var (
	ErrStreamRequestFailed = apperror.Unavailable("stream request failed")
)

type localLLMStreamingProvider struct {
	client          *httpclient.Client
	model           string
	temperature     float32
	topK            int32
	topP            float32
	maxOutputTokens int32
}

func New(
	baseURL string,
	apiKey string,
	model string,
	temperature float32,
	topK int32,
	topP float32,
	maxOutputTokens int32,
	timeout time.Duration,
) streamingchatPkg.Provider {
	clientBuilder := httpclient.NewClientBuilder().
		BaseURL(strings.TrimRight(baseURL, "/")).
		Timeout(timeout).
		JSONContentType()

	if apiKey != "" {
		clientBuilder.BearerAuth(apiKey)
	}

	return &localLLMStreamingProvider{
		client:          clientBuilder.Build(),
		model:           model,
		temperature:     temperature,
		topK:            topK,
		topP:            topP,
		maxOutputTokens: maxOutputTokens,
	}
}

func (p *localLLMStreamingProvider) StreamChat(ctx context.Context, systemInstruction string, messages []streamingchatPkg.Message) (<-chan streamingchatPkg.Chunk, error) {
	openAIMessages := make([]chatMessage, 0, len(messages)+1)
	if strings.TrimSpace(systemInstruction) != "" {
		openAIMessages = append(openAIMessages, chatMessage{
			Role:    "system",
			Content: systemInstruction,
		})
	}
	for _, msg := range messages {
		role := msg.Role
		if role == "model" {
			role = "assistant"
		}
		openAIMessages = append(openAIMessages, chatMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	payload := chatCompletionRequest{
		Model:       p.model,
		Messages:    openAIMessages,
		Temperature: p.temperature,
		TopK:        p.topK,
		TopP:        p.topP,
		MaxTokens:   p.maxOutputTokens,
		Stream:      true,
	}

	response, err := p.client.Post("/chat/completions").
		Context(ctx).
		Body(payload).
		DoStream()
	if err != nil {
		slog.Error("Local LLM streaming request failed",
			"component", "Generation",
			"error", err,
			"model", p.model,
		)
		return nil, ErrStreamRequestFailed
	}

	if !response.IsSuccess() {
		body := make([]byte, 4096)
		n, _ := response.BodyReader.Read(body)
		if err := response.BodyReader.Close(); err != nil {
			slog.Warn("Local LLM streaming body close failed",
				"component", "Generation",
				"error", err)
		}
		slog.Error("Local LLM streaming request returned non-success status",
			"component", "Generation",
			"status", response.Status,
			"body", string(body[:n]),
			"model", p.model,
		)
		return nil, ErrStreamRequestFailed
	}

	ch := make(chan streamingchatPkg.Chunk, 64)

	go func() {
		defer close(ch)
		defer func() {
			if err := response.BodyReader.Close(); err != nil {
				slog.Warn("Local LLM streaming body close failed",
					"component", "Generation",
					"error", err)
			}
		}()

		scanner := bufio.NewScanner(response.BodyReader)
		scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

		for scanner.Scan() {
			line := scanner.Text()

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if strings.TrimSpace(data) == "" {
				continue
			}

			if data == "[DONE]" {
				break
			}

			var sseResp SSEResponse
			if err := json.Unmarshal([]byte(data), &sseResp); err != nil {
				slog.Error("Failed to parse SSE response",
					"component", "Generation",
					"error", err,
					"data", data,
				)
				continue
			}

			if sseResp.Error != nil {
				slog.Error("Local LLM streaming returned error",
					"component", "Generation",
					"message", sseResp.Error.Message,
					"type", sseResp.Error.Type,
				)
				select {
				case ch <- streamingchatPkg.Chunk{
					Type:  streamingchatPkg.ChunkTypeError,
					Error: errors.New(sseResp.Error.Message),
				}:
				case <-ctx.Done():
				}
				return
			}

			for _, choice := range sseResp.Choices {
				if choice.Delta.Content == "" {
					continue
				}
				select {
				case ch <- streamingchatPkg.Chunk{
					Type:    streamingchatPkg.ChunkTypeText,
					Content: choice.Delta.Content,
				}:
				case <-ctx.Done():
					return
				}
			}
		}

		if err := scanner.Err(); err != nil {
			select {
			case ch <- streamingchatPkg.Chunk{
				Type:  streamingchatPkg.ChunkTypeError,
				Error: err,
			}:
			case <-ctx.Done():
			}
			return
		}

		select {
		case ch <- streamingchatPkg.Chunk{
			Type: streamingchatPkg.ChunkTypeDone,
		}:
		case <-ctx.Done():
		}
	}()

	return ch, nil
}
