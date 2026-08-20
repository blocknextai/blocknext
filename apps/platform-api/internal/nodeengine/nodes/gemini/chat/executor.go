package chat

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/gemini/helpers"
)

var (
	ErrExecutorEmptyResponse = apperror.Internal("empty response")
)

type GeminiChatExecutorInput struct {
	Model       string  `schema:"model"`
	Prompt      string  `schema:"prompt"`
	Temperature float64 `schema:"temperature"`
}

type GeminiChatExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[GeminiChatExecutorInput]
}

func NewGeminiChatExecutor(
	nodeID string,
	validator *jsonschema.Validator[GeminiChatExecutorInput],
) *GeminiChatExecutor {
	return &GeminiChatExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

type SuccessResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func (e *GeminiChatExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "gemini_api")
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
			response, err := client.Post("/models/"+input.Model+":generateContent").
				JSONContentType().
				Body(map[string]any{
					"contents": []map[string]any{
						{
							"parts": []map[string]any{
								{
									"text": input.Prompt,
								},
							},
						},
					},
					"generationConfig": map[string]any{
						"temperature": input.Temperature,
					},
				}).
				Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			if len(successResponse.Candidates) == 0 {
				return nil, ErrExecutorEmptyResponse
			}

			if len(successResponse.Candidates[0].Content.Parts) == 0 {
				return nil, ErrExecutorEmptyResponse
			}

			for _, candidate := range successResponse.Candidates {
				for _, part := range candidate.Content.Parts {
					results = append(results, map[string]any{
						"text": part.Text,
					})
				}
			}
		}

		return results, nil
	}
}
