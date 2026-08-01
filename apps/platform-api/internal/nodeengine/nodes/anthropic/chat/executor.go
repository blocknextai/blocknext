package chat

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/anthropic/helpers"
)

var (
	ErrExecutorEmptyResponse = apperror.Internal("empty response")
)

type AnthropicChatExecutorInput struct {
	Model       string  `schema:"model"`
	MaxTokens   int     `schema:"maxTokens"`
	Prompt      string  `schema:"prompt"`
	Temperature float64 `schema:"temperature"`
}

type AnthropicChatExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[AnthropicChatExecutorInput]
}

func NewAnthropicChatExecutor(
	nodeID string,
	validator *jsonschema.Validator[AnthropicChatExecutorInput],
) *AnthropicChatExecutor {
	return &AnthropicChatExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

type SuccessResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (e *AnthropicChatExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "anthropic_api")
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
			response, err := client.Post("/messages").
				JSONContentType().
				Body(map[string]any{
					"model":       input.Model,
					"max_tokens":  input.MaxTokens,
					"temperature": input.Temperature,
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
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			if len(successResponse.Content) == 0 {
				return nil, ErrExecutorEmptyResponse
			}

			for _, item := range successResponse.Content {
				results = append(results, map[string]any{
					"text": item.Text,
				})
			}
		}

		return results, nil
	}
}
