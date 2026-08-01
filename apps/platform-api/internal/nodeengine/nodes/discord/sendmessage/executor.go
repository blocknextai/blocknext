package sendmessage

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/discord/helpers"
)

var (
	ErrExecutorEmptyResponse = apperror.Internal("empty response")
)

type DiscordSendMessageExecutorInput struct {
	ChannelID string `schema:"channelId"`
	Content   string `schema:"content"`
}

type DiscordSendMessageExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[DiscordSendMessageExecutorInput]
}

func NewDiscordSendMessageExecutor(
	nodeID string,
	validator *jsonschema.Validator[DiscordSendMessageExecutorInput],
) *DiscordSendMessageExecutor {
	return &DiscordSendMessageExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

type SuccessResponse struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *DiscordSendMessageExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "discord_api")
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
			response, err := client.Post("/channels/"+input.ChannelID+"/messages").
				JSONContentType().
				Body(map[string]any{
					"content": input.Content,
				}).
				Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Message)
			}

			if successResponse.ID == "" {
				return nil, ErrExecutorEmptyResponse
			}

			results = append(results, map[string]any{
				"status":    true,
				"messageId": successResponse.ID,
			})
		}

		return results, nil
	}
}
