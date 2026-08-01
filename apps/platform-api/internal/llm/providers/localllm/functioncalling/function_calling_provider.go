package functioncalling

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/llm/functioncalling"
)

type localLLMProvider struct {
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
	baseURL string,
	apiKey string,
	model string,
	temperature float32,
	topK int32,
	topP float32,
	maxOutputTokens int32,
	timeout time.Duration,
	maxTimeout time.Duration,
) (functioncalling.FunctionCallingService, error) {
	builder := httpclient.NewClientBuilder().
		BaseURL(strings.TrimRight(baseURL, "/")).
		JSONContentType()
	if apiKey != "" {
		builder.BearerAuth(apiKey)
	}

	return &localLLMProvider{
		systemInstruction: systemInstruction,
		model:             model,
		temperature:       temperature,
		topK:              topK,
		topP:              topP,
		maxOutputTokens:   maxOutputTokens,
		timeout:           timeout,
		maxTimeout:        maxTimeout,
		client:            builder.Build(),
	}, nil
}

func (f *localLLMProvider) ExecuteWithContext(ctx context.Context, data []map[string]any, functionDeclarations []map[string]any) ([]map[string]any, error) {
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
	messages := make([]map[string]any, 0, len(input)+1)
	totalChars := 0
	if f.systemInstruction != "" {
		messages = append(messages, map[string]any{
			"role":    "system",
			"content": f.systemInstruction,
		})
		totalChars += len(f.systemInstruction)
	}

	for _, item := range input {
		builder.Reset()
		if item.Instruction != "" {
			builder.WriteString("NODE INSTRUCTION:\n")
			builder.WriteString(item.Instruction)
			builder.WriteString("\n")
		}
		if item.RuntimeInstruction != "" {
			builder.WriteString("RUNTIME INSTRUCTION:\n")
			builder.WriteString(item.RuntimeInstruction)
			builder.WriteString("\n")
		}
		if item.RuntimePrompt != "" {
			builder.WriteString("RUNTIME PROMPT (reference only):\n")
			builder.WriteString(item.RuntimePrompt)
			builder.WriteString("\n")
		}
		content := builder.String()
		totalChars += len(content)
		messages = append(messages, map[string]any{
			"role":    "user",
			"content": content,
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

	tools := make([]map[string]any, 0, len(functionDeclarations))
	for _, decl := range functionDeclarations {
		tools = append(tools, map[string]any{
			"type":     "function",
			"function": decl,
		})
	}

	requestBody := map[string]any{
		"model":       f.model,
		"messages":    messages,
		"temperature": temperature,
		"top_k":       f.topK,
		"top_p":       f.topP,
		"max_tokens":  f.maxOutputTokens,
		"tools":       tools,
		"tool_choice": "required",
	}

	reqCtx := ctx
	if dynamicTimeout > 0 {
		var cancel context.CancelFunc
		reqCtx, cancel = context.WithTimeout(ctx, dynamicTimeout)
		defer cancel()
	}

	var successResponse SuccessResponse
	var errorResponse ErrorResponse
	response, err := f.client.Post("/chat/completions").
		Context(reqCtx).
		Body(requestBody).
		Do(&successResponse, &errorResponse)
	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		if response.Status == 429 {
			return nil, functioncalling.ErrRateLimited
		}
		slog.Error("Local LLM function calling request failed",
			"component", "FunctionCalling",
			"status", response.Status,
			"upstream_message", errorResponse.Error.Message,
			"upstream_type", errorResponse.Error.Type,
			"model", f.model,
		)
		return nil, functioncalling.ErrProviderRequestFailed
	}

	if len(successResponse.Choices) == 0 {
		return nil, functioncalling.ErrNoCandidates
	}

	toolCalls := successResponse.Choices[0].Message.ToolCalls
	if len(toolCalls) == 0 {
		return nil, functioncalling.ErrNoContentParts
	}

	results := make([]map[string]any, 0, len(toolCalls))
	for _, call := range toolCalls {
		if strings.TrimSpace(call.Function.Name) == "" {
			return nil, functioncalling.ErrNoFunctionCallName
		}

		args := map[string]any{}
		if strings.TrimSpace(call.Function.Arguments) != "" {
			if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
				slog.Error("Failed to unmarshal tool call arguments",
					"component", "FunctionCalling",
					"function", call.Function.Name,
					"error", err)
				return nil, err
			}
		}

		results = append(results, args)
	}

	return results, nil
}
