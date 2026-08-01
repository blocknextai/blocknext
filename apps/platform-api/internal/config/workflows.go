package config

import (
	"time"
)

type WorkflowsOptions struct {
	Generation WorkflowsGenerationOptions `envPrefix:"GENERATION_"`
}

type WorkflowsGenerationProvider string

const (
	GeminiWorkflowsGenerationProvider   WorkflowsGenerationProvider = "gemini"
	LocalLLMWorkflowsGenerationProvider WorkflowsGenerationProvider = "local"
)

type WorkflowsGenerationOptions struct {
	Enabled               bool                               `env:"ENABLED"`
	Provider              WorkflowsGenerationProvider        `env:"PROVIDER"`
	SystemInstructionFile string                             `env:"SYSTEM_INSTRUCTION_FILE"`
	SystemInstruction     string                             `env:"-"`
	Gemini                WorkflowsGenerationGeminiOptions   `envPrefix:"GEMINI_"`
	LocalLLM              WorkflowsGenerationLocalLLMOptions `envPrefix:"LOCAL_"`
}

type WorkflowsGenerationGeminiOptions struct {
	APIKey          string        `env:"API_KEY"`
	Model           string        `env:"MODEL"`
	Timeout         time.Duration `env:"TIMEOUT"`
	Temperature     float32       `env:"TEMPERATURE"`
	TopK            int32         `env:"TOP_K"`
	TopP            float32       `env:"TOP_P"`
	MaxOutputTokens int32         `env:"MAX_OUTPUT_TOKENS"`
}

type WorkflowsGenerationLocalLLMOptions struct {
	BaseURL         string        `env:"BASE_URL"`
	APIKey          string        `env:"API_KEY"`
	Model           string        `env:"MODEL"`
	Timeout         time.Duration `env:"TIMEOUT"`
	Temperature     float32       `env:"TEMPERATURE"`
	TopK            int32         `env:"TOP_K"`
	TopP            float32       `env:"TOP_P"`
	MaxOutputTokens int32         `env:"MAX_OUTPUT_TOKENS"`
}
