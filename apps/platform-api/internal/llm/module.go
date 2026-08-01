package llm

import (
	"log/slog"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/platform-api/internal/config"
	"github.com/blocknextai/platform-api/internal/llm/functioncalling"
	geminiFunctionCalling "github.com/blocknextai/platform-api/internal/llm/providers/gemini/functioncalling"
	geminiStreamingChat "github.com/blocknextai/platform-api/internal/llm/providers/gemini/streamingchat"
	localLLMFunctionCalling "github.com/blocknextai/platform-api/internal/llm/providers/localllm/functioncalling"
	localLLMStreamingChat "github.com/blocknextai/platform-api/internal/llm/providers/localllm/streamingchat"
	"github.com/blocknextai/platform-api/internal/llm/streamingchat"
)

var (
	ErrInvalidStreamingChatProvider   = apperror.Internal("invalid streaming chat provider")
	ErrInvalidFunctionCallingProvider = apperror.Internal("invalid function calling provider")
)

type Dependencies struct {
	ChatProviderOptions            config.WorkflowsGenerationOptions
	FunctionCallingProviderOptions config.FunctionCallingOptions
}

type Module struct {
	StreamingChatProvider  streamingchat.Provider
	FunctionCallingService functioncalling.FunctionCallingService
}

func NewModule(deps Dependencies) (*Module, error) {
	streamingChatProvider, err := newStreamingChatProvider(deps.ChatProviderOptions)
	if err != nil {
		return nil, err
	}

	functionCallingService, err := newFunctionCallingService(deps.FunctionCallingProviderOptions)
	if err != nil {
		return nil, err
	}

	return &Module{
		StreamingChatProvider:  streamingChatProvider,
		FunctionCallingService: functionCallingService,
	}, nil
}

func newStreamingChatProvider(options config.WorkflowsGenerationOptions) (streamingchat.Provider, error) {
	if !options.Enabled {
		slog.Info("streaming chat provider is disabled")
		return nil, nil
	}

	if options.Provider == config.GeminiWorkflowsGenerationProvider {
		slog.Info("streaming chat provider initialized", "provider", options.Provider)
		return geminiStreamingChat.New(
			options.Gemini.APIKey,
			options.Gemini.Model,
			options.Gemini.Temperature,
			options.Gemini.TopK,
			options.Gemini.TopP,
			options.Gemini.MaxOutputTokens,
			options.Gemini.Timeout,
		), nil
	}

	if options.Provider == config.LocalLLMWorkflowsGenerationProvider {
		slog.Info("streaming chat provider initialized", "provider", options.Provider)
		return localLLMStreamingChat.New(
			options.LocalLLM.BaseURL,
			options.LocalLLM.APIKey,
			options.LocalLLM.Model,
			options.LocalLLM.Temperature,
			options.LocalLLM.TopK,
			options.LocalLLM.TopP,
			options.LocalLLM.MaxOutputTokens,
			options.LocalLLM.Timeout,
		), nil
	}

	return nil, ErrInvalidStreamingChatProvider
}

func newFunctionCallingService(options config.FunctionCallingOptions) (functioncalling.FunctionCallingService, error) {
	if !options.Enabled {
		slog.Info("function calling service is disabled")
		return nil, nil
	}

	if options.Provider == config.GeminiFunctionCallingProvider {
		slog.Info("function calling service initialized", "provider", options.Provider)
		return geminiFunctionCalling.New(
			options.SystemInstruction,
			options.Gemini.APIKey,
			options.Gemini.Model,
			options.Gemini.Temperature,
			options.Gemini.TopK,
			options.Gemini.TopP,
			options.Gemini.MaxOutputTokens,
			options.Gemini.Timeout,
			options.Gemini.MaxTimeout,
		)
	}

	if options.Provider == config.LocalLLMFunctionCallingProvider {
		slog.Info("function calling service initialized", "provider", options.Provider)
		return localLLMFunctionCalling.New(
			options.SystemInstruction,
			options.LocalLLM.BaseURL,
			options.LocalLLM.APIKey,
			options.LocalLLM.Model,
			options.LocalLLM.Temperature,
			options.LocalLLM.TopK,
			options.LocalLLM.TopP,
			options.LocalLLM.MaxOutputTokens,
			options.LocalLLM.Timeout,
			options.LocalLLM.MaxTimeout,
		)
	}

	return nil, ErrInvalidFunctionCallingProvider
}
