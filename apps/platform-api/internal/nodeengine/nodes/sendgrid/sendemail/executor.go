package sendemail

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/sendgrid/helpers"
)

var (
	ErrExecutorRequestFailed = apperror.Internal("request failed")
)

type SendgridSendEmailExecutorInput struct {
	From    string   `schema:"from"`
	To      []string `schema:"to"`
	Subject string   `schema:"subject"`
	Content string   `schema:"content"`
}

type SendgridSendEmailExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[SendgridSendEmailExecutorInput]
}

func NewSendgridSendEmailExecutor(
	nodeID string,
	validator *jsonschema.Validator[SendgridSendEmailExecutorInput],
) *SendgridSendEmailExecutor {
	return &SendgridSendEmailExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

type ErrorResponse struct {
	Errors []struct {
		Message string `json:"message"`
		Field   string `json:"field"`
	} `json:"errors"`
}

func (e *SendgridSendEmailExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "sendgrid_api")
		apiKey := credential.String("apiKey")
		client := helpers.CreateClient(ctx, apiKey)

		results := make([]map[string]any, 0, len(data))
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			toRecipients := make([]map[string]any, 0, len(input.To))
			for _, email := range input.To {
				toRecipients = append(toRecipients, map[string]any{
					"email": email,
				})
			}

			payload := map[string]any{
				"personalizations": []map[string]any{
					{
						"to": toRecipients,
					},
				},
				"from": map[string]any{
					"email": input.From,
				},
				"subject": input.Subject,
				"content": []map[string]any{
					{
						"type":  "text/html",
						"value": input.Content,
					},
				},
			}

			var errorResponse ErrorResponse
			response, err := client.Post("/mail/send").
				JSONContentType().
				Body(payload).
				Do(nil, &errorResponse)
			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				if len(errorResponse.Errors) > 0 {
					return nil, apperror.Internal(errorResponse.Errors[0].Message)
				}
				return nil, ErrExecutorRequestFailed
			}

			results = append(results, map[string]any{
				"status": true,
			})
		}

		return results, nil
	}
}
