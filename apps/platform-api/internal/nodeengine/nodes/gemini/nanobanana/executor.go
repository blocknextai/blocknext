package nanobanana

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/base64"
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/gemini/helpers"
)

var (
	ErrExecutorEmptyResponse = apperror.Internal("empty response")
	ErrExecutorEmptyPart     = apperror.Internal("empty part")
	ErrExecutorEmptyData     = apperror.Internal("empty data")
)

type GeminiImageGenerationExecutorInput struct {
	Prompt string `schema:"prompt"`
	Model  string `schema:"model"`
	Image  string `schema:"image"`
}

type GeminiImageGenerationExecutor struct {
	executors.Executor
	validator   *jsonschema.Validator[GeminiImageGenerationExecutorInput]
	fileGateway filegateway.FileGateway
}

func NewGeminiImageGenerationExecutor(
	nodeID string,
	validator *jsonschema.Validator[GeminiImageGenerationExecutorInput],
	fileGateway filegateway.FileGateway,
) *GeminiImageGenerationExecutor {
	return &GeminiImageGenerationExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator:   validator,
		fileGateway: fileGateway,
	}
}

type SuccessResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				InlineData *struct {
					Data string `json:"data"`
				} `json:"inlineData,omitempty"`
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

func (e *GeminiImageGenerationExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			parts, err := e.buildRequestParts(ctx, input)
			if err != nil {
				return nil, err
			}

			payload := e.buildRequestPayload(parts)

			var successResponse SuccessResponse
			var errorResponse ErrorResponse
			response, err := client.Post("/models/"+input.Model+":generateContent").
				JSONContentType().
				Body(payload).
				Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			imageURL, err := e.processResponse(ctx, successResponse)
			if err != nil {
				return nil, err
			}

			results = append(results, map[string]any{
				"image": imageURL,
			})
		}

		return results, nil
	}
}

func (e *GeminiImageGenerationExecutor) buildRequestParts(ctx context.Context, input *GeminiImageGenerationExecutorInput) ([]map[string]any, error) {
	parts := []map[string]any{
		{
			"text": input.Prompt,
		},
	}

	if input.Image != "" {
		file, err := e.fileGateway.DownloadFile(ctx, input.Image)
		if err != nil {
			return nil, err
		}
		defer func() {
			_ = file.DataReader.Close()
		}()

		base64Data, err := base64.EncodeReader(file.DataReader)
		if err != nil {
			return nil, err
		}

		parts = append(parts, map[string]any{
			"inline_data": map[string]any{
				"mime_type": file.ContentType,
				"data":      base64Data,
			},
		})
	}

	return parts, nil
}

func (e *GeminiImageGenerationExecutor) buildRequestPayload(parts []map[string]any) map[string]any {
	return map[string]any{
		"contents": []map[string]any{
			{
				"parts": parts,
			},
		},
		"generationConfig": map[string]any{
			"responseModalities": []string{
				"TEXT",
				"IMAGE",
			},
		},
	}
}

func (e *GeminiImageGenerationExecutor) processResponse(ctx context.Context, response SuccessResponse) (string, error) {
	if len(response.Candidates) == 0 {
		return "", ErrExecutorEmptyResponse
	}

	if len(response.Candidates[0].Content.Parts) == 0 {
		return "", ErrExecutorEmptyPart
	}

	var base64Data string
	for _, part := range response.Candidates[0].Content.Parts {
		if part.InlineData != nil && part.InlineData.Data != "" {
			base64Data = part.InlineData.Data
			break
		}
	}

	if base64Data == "" {
		return "", ErrExecutorEmptyData
	}

	reader := base64.DecodeReader(base64Data)

	uploadResult, err := e.fileGateway.UploadFile(
		ctx,
		"04ed7773-3a64-403b-a548-80a78d6ca11f",
		"gemini-image-generation.png",
		reader,
	)
	if err != nil {
		return "", err
	}

	return uploadResult.URL, nil
}
