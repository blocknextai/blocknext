package chat

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/chatgpt/helpers"
)

var (
	ErrExecutorEmptyResponse = apperror.Internal("empty response")
)

type ChatgptChatExecutorInput struct {
	Model       string  `schema:"model"`
	Prompt      string  `schema:"prompt"`
	Temperature float64 `schema:"temperature"`
}

type ChatgptChatExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[ChatgptChatExecutorInput]
}

func NewChatgptChatExecutor(
	nodeID string,
	validator *jsonschema.Validator[ChatgptChatExecutorInput],
) *ChatgptChatExecutor {
	return &ChatgptChatExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

type SuccessResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Param   string `json:"param"`
		Type    string `json:"type"`
	} `json:"error"`
}

func (e *ChatgptChatExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "chatgpt_api")
		apiKey := credential.String("apiKey")
		client := helpers.CreateClient(ctx, apiKey)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			var successResponse SuccessResponse
			var errorResponse ErrorResponse
			response, err := client.Post("/chat/completions").
				JSONContentType().
				Body(map[string]any{
					"model": input.Model,
					"messages": []map[string]any{
						{
							"role":    "user",
							"content": input.Prompt,
						},
					},
				}).
				Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Error.Code)
			}

			if len(successResponse.Choices) == 0 {
				return nil, ErrExecutorEmptyResponse
			}

			for _, item := range successResponse.Choices {
				results = append(results, map[string]any{
					"text": item.Message.Content,
				})
			}
		}
		return results, nil
	}
}
