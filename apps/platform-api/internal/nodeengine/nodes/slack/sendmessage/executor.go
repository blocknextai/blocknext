package sendmessage

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/slack/helpers"
)

var (
	ErrExecutorEmptyResponse = apperror.Internal("empty response")
)

type SlackSendMessageExecutorInput struct {
	Channel string `schema:"channel"`
	Text    string `schema:"text"`
}

type SlackSendMessageExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[SlackSendMessageExecutorInput]
}

func NewSlackSendMessageExecutor(
	nodeID string,
	validator *jsonschema.Validator[SlackSendMessageExecutorInput],
) *SlackSendMessageExecutor {
	return &SlackSendMessageExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

type SuccessResponse struct {
	Ok      bool   `json:"ok"`
	Channel string `json:"channel"`
	Ts      string `json:"ts"`
	Error   string `json:"error"`
}

func (e *SlackSendMessageExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "slack_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")
		client := helpers.CreateClient(ctx, accessToken)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			var successResponse SuccessResponse
			_, err = client.Post("/chat.postMessage").
				JSONContentType().
				Body(map[string]any{
					"channel": input.Channel,
					"text":    input.Text,
				}).
				Do(&successResponse, nil)

			if err != nil {
				return nil, err
			}

			if !successResponse.Ok {
				return nil, apperror.Internal(successResponse.Error)
			}

			if successResponse.Ts == "" {
				return nil, ErrExecutorEmptyResponse
			}

			results = append(results, map[string]any{
				"status":    true,
				"messageId": successResponse.Ts,
			})
		}

		return results, nil
	}
}
