package getpage

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

type NotionGetPageExecutorInput struct {
	PageID string `schema:"pageId"`
}

type NotionGetPageExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[NotionGetPageExecutorInput]
}

func NewNotionGetPageExecutor(
	nodeID string,
	validator *jsonschema.Validator[NotionGetPageExecutorInput],
) *NotionGetPageExecutor {
	return &NotionGetPageExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

func (e *NotionGetPageExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			getPageResponse, err := e.getNotionPage(client, *input)
			if err != nil {
				return nil, err
			}

			if getPageResponse == nil {
				return nil, ErrExecutorEmptyResponse
			}

			results = append(results, map[string]any{
				"id":         getPageResponse.ID,
				"object":     getPageResponse.Object,
				"url":        getPageResponse.URL,
				"properties": getPageResponse.Properties,
			})
		}

		return results, nil
	}
}

type GetPageResponse struct {
	ID         string         `json:"id"`
	Object     string         `json:"object"`
	URL        string         `json:"url"`
	Properties map[string]any `json:"properties"`
}

type ErrorResponse struct {
	Object  string `json:"object"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *NotionGetPageExecutor) getNotionPage(client *httpclient.Client, input NotionGetPageExecutorInput) (*GetPageResponse, error) {
	var successResponse GetPageResponse
	var errorResponse ErrorResponse
	response, err := client.Get("/pages/"+input.PageID).
		JSONContentType().
		Do(&successResponse, &errorResponse)

	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		return nil, apperror.Internal(errorResponse.Message)
	}

	return &successResponse, nil
}
