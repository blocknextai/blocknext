package streamingchat

import (
	"bufio"
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/go-packages/json"
	llmCache "github.com/blocknextai/platform-api/internal/llm/cache"
	"github.com/blocknextai/platform-api/internal/llm/providers/gemini"
	geminiCache "github.com/blocknextai/platform-api/internal/llm/providers/gemini/cache"
	streamingchatPkg "github.com/blocknextai/platform-api/internal/llm/streamingchat"
)

var (
	ErrStreamRequestFailed = apperror.Unavailable("stream request failed")
)

type geminiStreamingProvider struct {
	client          *httpclient.Client
	model           string
	temperature     float32
	topK            int32
	topP            float32
	maxOutputTokens int32
	cache           llmCache.Cache
}

func New(
	apiKey string,
	model string,
	temperature float32,
	topK int32,
	topP float32,
	maxOutputTokens int32,
	timeout time.Duration,
) streamingchatPkg.Provider {
	var pathBuilder strings.Builder
	pathBuilder.WriteString(gemini.BaseURL)
	pathBuilder.WriteString("/models/")
	pathBuilder.WriteString(model)
	pathBuilder.WriteString(":streamGenerateContent")

	client := httpclient.NewClientBuilder().
		BaseURL(pathBuilder.String()).
		QueryParam("alt", "sse").
		QueryParam("key", apiKey).
		Timeout(timeout).
		JSONContentType().
		Build()

	return &geminiStreamingProvider{
		client:          client,
		model:           model,
		temperature:     temperature,
		topK:            topK,
		topP:            topP,
		maxOutputTokens: maxOutputTokens,
		cache:           geminiCache.New(apiKey, model, "Generation", timeout),
	}
}

func (p *geminiStreamingProvider) StreamChat(ctx context.Context, systemInstruction string, messages []streamingchatPkg.Message) (<-chan streamingchatPkg.Chunk, error) {
	contents := make([]Content, 0, len(messages))
	for _, msg := range messages {
		contents = append(contents, Content{
			Role:  msg.Role,
			Parts: []Part{{Text: msg.Content}},
		})
	}

	payload := streamGenerateContentRequest{
		Contents: contents,
		GenerationConfig: GenerationConfig{
			Temperature:     p.temperature,
			TopK:            p.topK,
			TopP:            p.topP,
			MaxOutputTokens: p.maxOutputTokens,
		},
	}

	if strings.TrimSpace(systemInstruction) != "" {
		if cacheName := p.cache.Ensure(ctx, systemInstruction); cacheName != "" {
			payload.CachedContent = cacheName
		} else {
			payload.SystemInstruction = &SystemInstruction{
				Parts: []Part{{Text: systemInstruction}},
			}
		}
	}

	response, err := p.client.Post("").
		Context(ctx).
		Body(payload).
		DoStream()
	if err != nil {
		return nil, ErrStreamRequestFailed
	}

	if !response.IsSuccess() {
		body := make([]byte, 4096)
		n, _ := response.BodyReader.Read(body)
		if err := response.BodyReader.Close(); err != nil {
			slog.Warn("Gemini streaming body close failed",
				"component", "Generation",
				"error", err)
		}
		slog.Error("Gemini streaming request failed",
			"component", "Generation",
			"status", response.Status,
			"body", string(body[:n]),
		)
		return nil, ErrStreamRequestFailed
	}

	ch := make(chan streamingchatPkg.Chunk, 64)
	cachedContentUsed := payload.CachedContent != ""

	go func() {
		var usage *UsageMetadata
		defer func() {
			if usage == nil {
				return
			}
			slog.Info("Gemini streaming token usage",
				"component", "Generation",
				"model", p.model,
				"cached_content_used", cachedContentUsed,
				"prompt_tokens", usage.PromptTokenCount,
				"cached_tokens", usage.CachedContentTokenCount,
				"candidates_tokens", usage.CandidatesTokenCount,
				"thoughts_tokens", usage.ThoughtsTokenCount,
				"total_tokens", usage.TotalTokenCount,
			)
		}()

		defer close(ch)
		defer func() {
			if err := response.BodyReader.Close(); err != nil {
				slog.Warn("Gemini streaming body close failed",
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

			var sseResp SSEResponse
			if err := json.Unmarshal([]byte(data), &sseResp); err != nil {
				slog.Error("Failed to parse SSE response",
					"component", "Generation",
					"error", err,
					"data", data,
				)
				continue
			}

			if sseResp.UsageMetadata != nil {
				usage = sseResp.UsageMetadata
			}

			for _, candidate := range sseResp.Candidates {
				for _, p := range candidate.Content.Parts {
					if p.Text != "" {
						select {
						case ch <- streamingchatPkg.Chunk{
							Type:    streamingchatPkg.ChunkTypeText,
							Content: p.Text,
						}:
						case <-ctx.Done():
							return
						}
					}
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
