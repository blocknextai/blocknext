package createorupdate

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/docs/helpers"
)

type GoogleDocsCreateOrUpdateExecutorInput struct {
	Title      string `schema:"title"`
	Content    string `schema:"content"`
	DocumentID string `schema:"documentId"`
}

type GoogleDocsCreateOrUpdateExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[GoogleDocsCreateOrUpdateExecutorInput]
}

func NewGoogleDocsCreateOrUpdateExecutor(
	nodeID string,
	validator *jsonschema.Validator[GoogleDocsCreateOrUpdateExecutorInput],
) *GoogleDocsCreateOrUpdateExecutor {
	return &GoogleDocsCreateOrUpdateExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

type GoogleDocsCreateOrUpdateResponse struct {
	DocumentID string `json:"documentId"`
	Title      string `json:"title"`
}

type GoogleDocsResponse struct {
	DocumentID string `json:"documentId"`
	Body       struct {
		Content []struct {
			EndIndex int `json:"endIndex"`
		} `json:"content"`
	} `json:"body"`
}

type GoogleDocsCreateOrUpdateErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func (e *GoogleDocsCreateOrUpdateExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "google_docs_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")

		client := helpers.CreateClient(ctx, accessToken)
		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			endIndex := 1
			documentID := input.DocumentID
			content := input.Content
			if documentID != "" {
				var docResponse GoogleDocsResponse
				var docError GoogleDocsCreateOrUpdateErrorResponse
				docRes, err := client.Get("/documents/"+documentID).
					Do(&docResponse, &docError)

				if err != nil {
					return nil, err
				}

				if !docRes.IsSuccess() {
					return nil, apperror.Internal(docError.Error.Message)
				}

				if len(docResponse.Body.Content) > 0 {
					endIndex = docResponse.Body.Content[len(docResponse.Body.Content)-1].EndIndex - 1
					content = "\n" + content
				}
			} else {
				var createResponse GoogleDocsCreateOrUpdateResponse
				var createError GoogleDocsCreateOrUpdateErrorResponse
				createRes, err := client.Post("/documents").
					JSONContentType().
					Body(map[string]any{
						"title": input.Title,
					}).
					Do(&createResponse, &createError)

				if err != nil {
					return nil, err
				}

				if !createRes.IsSuccess() {
					return nil, apperror.Internal(createError.Error.Message)
				}

				documentID = createResponse.DocumentID
			}

			var updateError GoogleDocsCreateOrUpdateErrorResponse
			response, err := client.Post("/documents/"+documentID+":batchUpdate").
				JSONContentType().
				Body(map[string]any{
					"requests": []map[string]any{
						{
							"insertText": map[string]any{
								"location": map[string]any{
									"index": endIndex,
								},
								"text": content,
							},
						},
					},
				}).
				Do(nil, &updateError)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(updateError.Error.Message)
			}

			results = append(results, map[string]any{
				"status": true,
			})
		}

		return results, nil
	}
}
