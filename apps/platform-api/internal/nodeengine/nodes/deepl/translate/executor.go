package translate

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/deepl/helpers"
)

var (
	ErrExecutorEmptyResponse = apperror.Internal("empty response")
)

type DeeplTranslateExecutorInput struct {
	Text           string `schema:"text"`
	SourceLanguage string `schema:"sourceLanguage"`
	TargetLanguage string `schema:"targetLanguage"`
}

type DeeplTranslateExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[DeeplTranslateExecutorInput]
}

func NewDeeplTranslateExecutor(
	nodeID string,
	validator *jsonschema.Validator[DeeplTranslateExecutorInput],
) *DeeplTranslateExecutor {
	return &DeeplTranslateExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

type TranslationResponse struct {
	Translations []struct {
		Text                   string `json:"text"`
		DetectedSourceLanguage string `json:"detected_source_language"`
	} `json:"translations"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (e *DeeplTranslateExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "deepl_api")
		apiKey := credential.String("apiKey")
		client := helpers.CreateClient(ctx, apiKey)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			requestBody := map[string]any{
				"text": []string{
					input.Text,
				},
				"target_lang": input.TargetLanguage,
			}

			if input.SourceLanguage != "auto" {
				requestBody["source_lang"] = input.SourceLanguage
			}

			var successResponse TranslationResponse
			var errorResponse ErrorResponse
			response, err := client.Post("/translate").
				JSONContentType().
				Body(requestBody).
				Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Message)
			}

			if len(successResponse.Translations) == 0 {
				return nil, ErrExecutorEmptyResponse
			}

			for _, translation := range successResponse.Translations {
				results = append(results, map[string]any{
					"text": translation.Text,
				})
			}
		}

		return results, nil
	}
}
