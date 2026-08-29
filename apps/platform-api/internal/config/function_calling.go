package config

import (
	"time"
)

type FunctionCallingProvider string

const (
	GeminiFunctionCallingProvider   FunctionCallingProvider = "gemini"
	LocalLLMFunctionCallingProvider FunctionCallingProvider = "local"
)

type FunctionCallingOptions struct {
	Enabled           bool                    `env:"ENABLED"`
	Provider          FunctionCallingProvider `env:"PROVIDER"`
	SystemInstruction string                  `env:"-"`

	Gemini   GeminiFunctionCallingOptions   `envPrefix:"GEMINI_"`
	LocalLLM LocalLLMFunctionCallingOptions `envPrefix:"LOCAL_"`
}

type GeminiFunctionCallingOptions struct {
	APIKey          string        `env:"API_KEY"`
	Model           string        `env:"MODEL"`
	Timeout         time.Duration `env:"TIMEOUT"`
	MaxTimeout      time.Duration `env:"MAX_TIMEOUT"`
	Temperature     float32       `env:"TEMPERATURE"`
	TopK            int32         `env:"TOP_K"`
	TopP            float32       `env:"TOP_P"`
	MaxOutputTokens int32         `env:"MAX_OUTPUT_TOKENS"`
}

type LocalLLMFunctionCallingOptions struct {
	BaseURL string `env:"BASE_URL"`
	APIKey  string `env:"API_KEY"`
	Model   string `env:"MODEL"`

	Timeout    time.Duration `env:"TIMEOUT"`
	MaxTimeout time.Duration `env:"MAX_TIMEOUT"`

	Temperature     float32 `env:"TEMPERATURE"`
	TopK            int32   `env:"TOP_K"`
	TopP            float32 `env:"TOP_P"`
	MaxOutputTokens int32   `env:"MAX_OUTPUT_TOKENS"`
}
