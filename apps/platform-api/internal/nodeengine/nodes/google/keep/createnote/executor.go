package createnote

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/keep/helpers"
)

type GoogleKeepCreateNoteExecutorInput struct {
	Title   string `schema:"title"`
	Content string `schema:"content"`
}

type GoogleKeepCreateNoteExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[GoogleKeepCreateNoteExecutorInput]
}

func NewGoogleKeepCreateNoteExecutor(
	nodeID string,
	validator *jsonschema.Validator[GoogleKeepCreateNoteExecutorInput],
) *GoogleKeepCreateNoteExecutor {
	return &GoogleKeepCreateNoteExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

func (e *GoogleKeepCreateNoteExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "google_keep_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")
		client := helpers.CreateClient(ctx, accessToken)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			var successResponse helpers.Note
			var errorResponse helpers.ErrorResponse
			response, err := client.Post("/notes").
				JSONContentType().
				Body(map[string]any{
					"title": input.Title,
					"body": helpers.Section{
						Text: &helpers.TextContent{
							Text: input.Content,
						},
					},
				}).
				Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			results = append(results, successResponse.ToMap())
		}

		return results, nil
	}
}
