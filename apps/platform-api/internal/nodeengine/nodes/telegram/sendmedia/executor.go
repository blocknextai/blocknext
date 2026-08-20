package sendmedia

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/telegram/helpers"
)

var (
	errMediaSendFailed = apperror.Internal("media send failed")
)

type TelegramSendMediaExecutorInput struct {
	ChatID    string   `schema:"chatId"`
	MediaURLs []string `schema:"mediaUrls"`
}

type TelegramSendMediaExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[TelegramSendMediaExecutorInput]
}

func NewTelegramSendMediaExecutor(
	nodeID string,
	validator *jsonschema.Validator[TelegramSendMediaExecutorInput],
) *TelegramSendMediaExecutor {
	return &TelegramSendMediaExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

type SuccessResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageID int64 `json:"message_id"`
	} `json:"result"`
}

type ErrorResponse struct {
	Description string `json:"description"`
	ErrorCode   string `json:"error_code"`
	Ok          bool   `json:"ok"`
}

func detectMediaType(url string) string {
	ext := strings.ToLower(filepath.Ext(url))

	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp":
		return "photo"
	case ".mp4", ".avi", ".mov", ".mkv", ".webm", ".flv", ".m4v":
		return "video"
	case ".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a", ".wma":
		return "audio"
	default:
		return "document"
	}
}

func (e *TelegramSendMediaExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "telegram_api")
		botToken := credential.String("botToken")
		client := helpers.CreateClient(ctx, botToken)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			for _, mediaURL := range input.MediaURLs {
				messageID, err := e.sendSingleMediaByURL(client, input.ChatID, mediaURL)
				if err != nil {
					results = append(results, map[string]any{
						"status": false,
						"error":  err.Error(),
					})
					continue
				}
				results = append(results, map[string]any{
					"status":    true,
					"messageId": messageID,
				})
			}
		}

		return results, nil
	}
}

func (e *TelegramSendMediaExecutor) sendSingleMediaByURL(client *httpclient.Client, chatID, mediaURL string) (int64, error) {
	mediaType := detectMediaType(mediaURL)

	var endpoint string
	var mediaField string
	switch mediaType {
	case "photo":
		endpoint = "/sendPhoto"
		mediaField = "photo"
	case "video":
		endpoint = "/sendVideo"
		mediaField = "video"
	case "audio":
		endpoint = "/sendAudio"
		mediaField = "audio"
	case "document":
		endpoint = "/sendDocument"
		mediaField = "document"
	default:
		endpoint = "/sendPhoto"
		mediaField = "photo"
	}

	var successResponse SuccessResponse
	var errorResponse ErrorResponse
	response, err := client.Post(endpoint).
		JSONContentType().
		Body(map[string]any{
			"chat_id":  chatID,
			mediaField: mediaURL,
		}).
		Do(&successResponse, &errorResponse)

	if err != nil {
		return 0, err
	}

	if !response.IsSuccess() {
		return 0, apperror.Internal(errorResponse.Description)
	}

	if !successResponse.Ok {
		return 0, errMediaSendFailed
	}

	return successResponse.Result.MessageID, nil
}
