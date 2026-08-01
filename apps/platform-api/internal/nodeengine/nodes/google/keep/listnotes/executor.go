package listnotes

import (
	"context"
	"strconv"
	"strings"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/keep/helpers"
)

type GoogleKeepListNotesExecutorInput struct {
	PageSize  int    `schema:"pageSize"`
	PageToken string `schema:"pageToken"`
	Filter    string `schema:"filter"`
}

type GoogleKeepListNotesExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[GoogleKeepListNotesExecutorInput]
}

func NewGoogleKeepListNotesExecutor(
	nodeID string,
	validator *jsonschema.Validator[GoogleKeepListNotesExecutorInput],
) *GoogleKeepListNotesExecutor {
	return &GoogleKeepListNotesExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

func (e *GoogleKeepListNotesExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			request := client.Get("/notes")
			if input.PageSize > 0 {
				request = request.QueryParam("pageSize", strconv.Itoa(input.PageSize))
			}
			if strings.TrimSpace(input.PageToken) != "" {
				request = request.QueryParam("pageToken", input.PageToken)
			}
			if strings.TrimSpace(input.Filter) != "" {
				request = request.QueryParam("filter", input.Filter)
			}

			var successResponse helpers.ListNotesResponse
			var errorResponse helpers.ErrorResponse
			response, err := request.Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			notes := make([]map[string]any, 0, len(successResponse.Notes))
			for _, note := range successResponse.Notes {
				notes = append(notes, note.ToMap())
			}

			results = append(results, map[string]any{
				"notes":         notes,
				"nextPageToken": successResponse.NextPageToken,
			})
		}

		return results, nil
	}
}
