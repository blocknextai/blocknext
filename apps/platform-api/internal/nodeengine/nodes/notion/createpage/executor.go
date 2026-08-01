package createpage

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/notion/helpers"
)

var (
	ErrExecutorEmptyResponse = apperror.Internal("empty response")
)

type NotionCreatePageExecutorInput struct {
	ParentID string `schema:"parentId"`
	Title    string `schema:"title"`
	Content  string `schema:"content"`
}

type NotionCreatePageExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[NotionCreatePageExecutorInput]
}

func NewNotionCreatePageExecutor(
	nodeID string,
	validator *jsonschema.Validator[NotionCreatePageExecutorInput],
) *NotionCreatePageExecutor {
	return &NotionCreatePageExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

func (e *NotionCreatePageExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "notion_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")
		client := helpers.CreateClient(ctx, accessToken)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			createPageResponse, err := e.createNotionPage(client, *input)
			if err != nil {
				return nil, err
			}

			if createPageResponse == nil {
				return nil, ErrExecutorEmptyResponse
			}

			results = append(results, map[string]any{
				"id": createPageResponse.ID,
			})
		}

		return results, nil
	}
}

type CreatePageResponse struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Object string `json:"object"`
}

type ErrorResponse struct {
	Object  string `json:"object"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *NotionCreatePageExecutor) createNotionPage(client *httpclient.Client, input NotionCreatePageExecutorInput) (*CreatePageResponse, error) {
	payload := map[string]any{
		"parent": map[string]any{
			"page_id": input.ParentID,
		},
		"properties": map[string]any{
			"Name": map[string]any{
				"title": []map[string]any{
					{
						"text": map[string]any{
							"content": input.Title,
						},
					},
				},
			},
		},
	}

	if input.Content != "" {
		payload["children"] = []map[string]any{
			{
				"object": "block",
				"type":   "paragraph",
				"paragraph": map[string]any{
					"rich_text": []map[string]any{
						{
							"type": "text",
							"text": map[string]any{
								"content": input.Content,
							},
						},
					},
				},
			},
		}
	}

	var successResponse CreatePageResponse
	var errorResponse ErrorResponse
	response, err := client.Post("/pages").
		JSONContentType().
		Body(payload).
		Do(&successResponse, &errorResponse)

	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		return nil, apperror.Internal(errorResponse.Message)
	}

	return &successResponse, nil
}
