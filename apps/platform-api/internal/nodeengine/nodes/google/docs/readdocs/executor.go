package readdocs

import (
	"context"
	"strings"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/docs/helpers"
)

type GoogleDocsReadDocsExecutorInput struct {
	DocumentID string `schema:"documentId"`
}

type GoogleDocsReadDocsExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[GoogleDocsReadDocsExecutorInput]
}

func NewGoogleDocsReadDocsExecutor(
	nodeID string,
	validator *jsonschema.Validator[GoogleDocsReadDocsExecutorInput],
) *GoogleDocsReadDocsExecutor {
	return &GoogleDocsReadDocsExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

type GoogleDocsReadDocsResponse struct {
	DocumentID string `json:"documentId"`
	Title      string `json:"title"`
	Body       struct {
		Content []struct {
			Paragraph struct {
				Elements []struct {
					TextRun struct {
						Content string `json:"content"`
					} `json:"textRun"`
				} `json:"elements"`
			} `json:"paragraph"`
		} `json:"content"`
	} `json:"body"`
}

type GoogleDocsReadDocsErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func (e *GoogleDocsReadDocsExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			var response GoogleDocsReadDocsResponse
			var errorResponse GoogleDocsReadDocsErrorResponse

			res, err := client.Get("/documents/"+input.DocumentID).
				Do(&response, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !res.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			var content strings.Builder
			for _, contentItem := range response.Body.Content {
				if contentItem.Paragraph.Elements != nil {
					for _, element := range contentItem.Paragraph.Elements {
						if element.TextRun.Content != "" {
							content.WriteString(element.TextRun.Content)
						}
					}
				}
			}

			results = append(results, map[string]any{
				"text": content.String(),
			})
		}

		return results, nil
	}
}
