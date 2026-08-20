package searchpages

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

type NotionSearchPagesExecutorInput struct {
	Query string  `schema:"query"`
	Limit float64 `schema:"limit"`
}

type NotionSearchPagesExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[NotionSearchPagesExecutorInput]
}

func NewNotionSearchPagesExecutor(
	nodeID string,
	validator *jsonschema.Validator[NotionSearchPagesExecutorInput],
) *NotionSearchPagesExecutor {
	return &NotionSearchPagesExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

func (e *NotionSearchPagesExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			searchResponse, err := e.searchNotionPages(client, *input)
			if err != nil {
				return nil, err
			}

			if searchResponse == nil {
				return nil, ErrExecutorEmptyResponse
			}

			results = append(results, map[string]any{
				"results": searchResponse.Results,
			})
		}

		return results, nil
	}
}

type SearchResponse struct {
	Object     string           `json:"object"`
	Results    []map[string]any `json:"results"`
	NextCursor string           `json:"next_cursor"`
	HasMore    bool             `json:"has_more"`
}

type ErrorResponse struct {
	Object  string `json:"object"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *NotionSearchPagesExecutor) searchNotionPages(client *httpclient.Client, input NotionSearchPagesExecutorInput) (*SearchResponse, error) {
	payload := map[string]any{
		"query": input.Query,
		"filter": map[string]any{
			"value":    "page",
			"property": "object",
		},
		"page_size": int(input.Limit),
	}

	var successResponse SearchResponse
	var errorResponse ErrorResponse
	response, err := client.Post("/search").
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
