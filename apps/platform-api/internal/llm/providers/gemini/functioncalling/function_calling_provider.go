package functioncalling

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/llm/functioncalling"
	"github.com/blocknextai/platform-api/internal/llm/providers/gemini"
)

type geminiProvider struct {
	systemInstruction string
	model             string
	temperature       float32
	topK              int32
	topP              float32
	maxOutputTokens   int32
	timeout           time.Duration
	maxTimeout        time.Duration
	client            *httpclient.Client
}

func New(
	systemInstruction string,
	apiKey string,
	model string,
	temperature float32,
	topK int32,
	topP float32,
	maxOutputTokens int32,
	timeout time.Duration,
	maxTimeout time.Duration,
) (functioncalling.FunctionCallingService, error) {
	var pathBuilder strings.Builder
	pathBuilder.WriteString(gemini.BaseURL)
	pathBuilder.WriteString("/models/")
	pathBuilder.WriteString(model)
	pathBuilder.WriteString(":generateContent")

	client := httpclient.NewClientBuilder().
		BaseURL(pathBuilder.String()).
		QueryParam("key", apiKey).
		JSONContentType().
		Build()

	return &geminiProvider{
		systemInstruction: systemInstruction,
		model:             model,
		temperature:       temperature,
		topK:              topK,
		topP:              topP,
		maxOutputTokens:   maxOutputTokens,
		timeout:           timeout,
		maxTimeout:        maxTimeout,
		client:            client,
	}, nil
}

func (f *geminiProvider) ExecuteWithContext(ctx context.Context, data []map[string]any, functionDeclarations []map[string]any) ([]map[string]any, error) {
	if len(functionDeclarations) == 0 {
		return nil, functioncalling.ErrEmptyFunctionDeclarations
	}

	if data == nil {
		return nil, functioncalling.ErrEmptyData
	}

	var input []functioncalling.FunctionCallingInput
	if err := json.ArgsToStruct(data, &input); err != nil {
		slog.Error("Failed to convert data to struct",
			"component", "FunctionCalling",
			"error", err)
		return nil, err
	}

	var builder strings.Builder
	contents := make([]map[string]any, 0, len(input))
	totalChars := 0
	for _, item := range input {
		userParts := make([]map[string]any, 0, 3)
		if item.Instruction != "" {
			builder.Reset()
			builder.WriteString("NODE INSTRUCTION:\n")
			builder.WriteString(item.Instruction)
			builder.WriteString("\n")
			text := builder.String()
			totalChars += len(text)
			userParts = append(userParts, map[string]any{"text": text})
		}
		if item.RuntimeInstruction != "" {
			builder.Reset()
			builder.WriteString("RUNTIME INSTRUCTION:\n")
			builder.WriteString(item.RuntimeInstruction)
			builder.WriteString("\n")
			text := builder.String()
			totalChars += len(text)
			userParts = append(userParts, map[string]any{"text": text})
		}
		if item.RuntimePrompt != "" {
			builder.Reset()
			builder.WriteString("RUNTIME PROMPT (reference only):\n")
			builder.WriteString(item.RuntimePrompt)
			builder.WriteString("\n")
			text := builder.String()
			totalChars += len(text)
			userParts = append(userParts, map[string]any{"text": text})
		}

		contents = append(contents, map[string]any{
			"role":  "user",
			"parts": userParts,
		})
	}

	dynamicTimeout := f.timeout + time.Duration(totalChars/1000)*time.Second
	if f.maxTimeout > 0 && dynamicTimeout > f.maxTimeout {
		dynamicTimeout = f.maxTimeout
	}

	temperature := f.temperature
	if temperature < 0.01 {
		temperature = 0.0
	}

	reqCtx := ctx
	if dynamicTimeout > 0 {
		var cancel context.CancelFunc
		reqCtx, cancel = context.WithTimeout(ctx, dynamicTimeout)
		defer cancel()
	}

	requestBody := map[string]any{
		"contents": contents,
		"systemInstruction": map[string]any{
			"parts": []map[string]any{
				{"text": f.systemInstruction},
			},
		},
		"generationConfig": map[string]any{
			"temperature":      temperature,
			"topK":             f.topK,
			"topP":             f.topP,
			"maxOutputTokens":  f.maxOutputTokens,
			"responseMimeType": "text/plain",
			"candidateCount":   1,
		},
		"tools": []map[string]any{
			{
				"functionDeclarations": functionDeclarations,
			},
		},
		"toolConfig": map[string]any{
			"functionCallingConfig": map[string]any{
				"mode":                 "ANY",
				"allowedFunctionNames": extractFunctionNames(functionDeclarations),
			},
		},
	}

	var successResponse SuccessResponse
	var errorResponse gemini.ErrorResponse
	response, err := f.client.Post("").
		Context(reqCtx).
		Body(requestBody).
		Do(&successResponse, &errorResponse)

	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		if errorResponse.Error.Code == 429 {
			return nil, functioncalling.ErrRateLimited
		}
		slog.Error("Gemini function calling request failed",
			"component", "FunctionCalling",
			"status", response.Status,
			"upstream_code", errorResponse.Error.Code,
			"upstream_message", errorResponse.Error.Message,
			"upstream_status", errorResponse.Error.Status,
			"model", f.model,
		)
		return nil, functioncalling.ErrProviderRequestFailed
	}

	usage := successResponse.UsageMetadata
	slog.Info("Gemini function calling token usage",
		"component", "FunctionCalling",
		"model", f.model,
		"prompt_tokens", usage.PromptTokenCount,
		"cached_tokens", usage.CachedContentTokenCount,
		"candidates_tokens", usage.CandidatesTokenCount,
		"thoughts_tokens", usage.ThoughtsTokenCount,
		"total_tokens", usage.TotalTokenCount,
	)

	if len(successResponse.Candidates) == 0 {
		return nil, functioncalling.ErrNoCandidates
	}

	parts := successResponse.Candidates[0].Content.Parts
	if len(parts) == 0 {
		return nil, functioncalling.ErrNoContentParts
	}

	results := make([]map[string]any, 0, len(parts))
	for _, part := range parts {
		functionCall := part.FunctionCall
		if strings.TrimSpace(functionCall.Name) == "" {
			return nil, functioncalling.ErrNoFunctionCallName
		}

		results = append(results, functionCall.Args)
	}

	return results, nil
}

func extractFunctionNames(declarations []map[string]any) []string {
	names := make([]string, 0, len(declarations))
	for _, decl := range declarations {
		if name, ok := decl["name"].(string); ok {
			names = append(names, name)
		}
	}
	return names
}
