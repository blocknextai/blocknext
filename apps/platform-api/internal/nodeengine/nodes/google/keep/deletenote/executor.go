package deletenote

import (
	"context"
	"strings"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/keep/helpers"
)

type GoogleKeepDeleteNoteExecutorInput struct {
	Name string `schema:"name"`
}

type GoogleKeepDeleteNoteExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[GoogleKeepDeleteNoteExecutorInput]
}

func NewGoogleKeepDeleteNoteExecutor(
	nodeID string,
	validator *jsonschema.Validator[GoogleKeepDeleteNoteExecutorInput],
) *GoogleKeepDeleteNoteExecutor {
	return &GoogleKeepDeleteNoteExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

func (e *GoogleKeepDeleteNoteExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			noteName := strings.TrimPrefix(input.Name, "/")

			var errorResponse helpers.ErrorResponse
			response, err := client.Delete("/"+noteName).
				Do(nil, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			results = append(results, map[string]any{
				"status": true,
			})
		}

		return results, nil
	}
}
