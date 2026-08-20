package updatepage

import (
	"context"
	"maps"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/notion/helpers"
)

var (
	ErrExecutorEmptyResponse   = apperror.Internal("empty response")
	ErrInvalidPropertiesFormat = apperror.Internal("invalid properties format")
)

type NotionUpdatePageExecutorInput struct {
	PageID     string `schema:"pageId"`
	Title      string `schema:"title"`
	Properties string `schema:"properties"`
}

type NotionUpdatePageExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[NotionUpdatePageExecutorInput]
}

func NewNotionUpdatePageExecutor(
	nodeID string,
	validator *jsonschema.Validator[NotionUpdatePageExecutorInput],
) *NotionUpdatePageExecutor {
	return &NotionUpdatePageExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

func (e *NotionUpdatePageExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			updateResponse, err := e.updateNotionPage(client, *input)
			if err != nil {
				return nil, err
			}

			if updateResponse == nil {
				return nil, ErrExecutorEmptyResponse
			}

			results = append(results, map[string]any{
				"status": true,
			})
		}

		return results, nil
	}
}

type UpdatePageResponse struct {
	ID         string         `json:"id"`
	Object     string         `json:"object"`
	Properties map[string]any `json:"properties"`
}

type ErrorResponse struct {
	Object  string `json:"object"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *NotionUpdatePageExecutor) updateNotionPage(client *httpclient.Client, input NotionUpdatePageExecutorInput) (*UpdatePageResponse, error) {
	payload := map[string]any{}
	if input.Title != "" {
		payload["properties"] = map[string]any{
			"Name": map[string]any{
				"title": []map[string]any{
					{
						"text": map[string]any{
							"content": input.Title,
						},
					},
				},
			},
		}
	}

	if input.Properties != "" {
		var additionalProperties map[string]any
		err := json.Unmarshal([]byte(input.Properties), &additionalProperties)
		if err != nil {
			return nil, ErrInvalidPropertiesFormat
		}

		if payload["properties"] == nil {
			payload["properties"] = additionalProperties
		} else {
			existingProps, ok := payload["properties"].(map[string]any)
			if !ok {
				return nil, ErrInvalidPropertiesFormat
			}
			maps.Copy(existingProps, additionalProperties)
		}
	}

	var successResponse UpdatePageResponse
	var errorResponse ErrorResponse
	response, err := client.Patch("/pages/"+input.PageID).
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
