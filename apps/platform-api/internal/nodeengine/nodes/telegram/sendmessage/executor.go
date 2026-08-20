package sendmessage

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/telegram/helpers"
)

var (
	errFailedToSendMessage = apperror.Internal("failed to send message")
	errMessageSendFailed   = apperror.Internal("message send failed")
)

type TelegramSendMessageExecutorInput struct {
	ChatID string `schema:"chatId"`
	Text   string `schema:"text"`
}

type TelegramSendMessageExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[TelegramSendMessageExecutorInput]
}

func NewTelegramSendMessageExecutor(
	nodeID string,
	validator *jsonschema.Validator[TelegramSendMessageExecutorInput],
) *TelegramSendMessageExecutor {
	return &TelegramSendMessageExecutor{
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

func (e *TelegramSendMessageExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			var successResponse SuccessResponse
			var errorResponse ErrorResponse
			response, err := client.Post("/sendMessage").
				JSONContentType().
				Body(map[string]any{
					"chat_id": input.ChatID,
					"text":    input.Text,
				}).
				Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, errFailedToSendMessage
			}

			if !successResponse.Ok {
				return nil, errMessageSendFailed
			}

			results = append(results, map[string]any{
				"status":    true,
				"messageId": successResponse.Result.MessageID,
			})
		}

		return results, nil
	}
}
