package config

import (
	"os"

	"github.com/caarlos0/env/v11"
)

type SharedConfig struct {
	AppEnv          AppEnv                 `env:"APP_ENV"`
	Database        DatabaseOptions        `envPrefix:"DATABASE_"`
	Cache           CacheOptions           `envPrefix:"CACHE_"`
	Broker          BrokerOptions          `envPrefix:"BROKER_"`
	EventBus        EventBusOptions        `envPrefix:"EVENT_BUS_"`
	SecretManager   SecretManagerOptions   `envPrefix:"SECRET_MANAGER_"`
	FileGateway     FileGatewayOptions     `envPrefix:"FILE_GATEWAY_"`
	JWT             JWTOptions             `envPrefix:"JWT_"`
	PlatformUI      PlatformUIOptions      `envPrefix:"PLATFORM_UI_"`
	Platform        PlatformOptions        `envPrefix:"PLATFORM_"`
	Webhook         WebhookOptions         `envPrefix:"WEBHOOK_"`
	Workflows       WorkflowsOptions       `envPrefix:"WORKFLOWS_"`
	FunctionCalling FunctionCallingOptions `envPrefix:"FUNCTION_CALLING_"`
	CredentialOAuth CredentialOAuthOptions `envPrefix:"CREDENTIAL_OAUTH_"`
	MCP             MCPOptions             `envPrefix:"MCP_"`
	Auth            AuthOptions            `envPrefix:"AUTH_"`
	EmailSender     EmailSenderOptions     `envPrefix:"EMAIL_SENDER_"`
	Semaphore       SemaphoreOptions       `envPrefix:"SEMAPHORE_"`
}

func LoadShared() (*SharedConfig, error) {
	cfg := &SharedConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	if cfg.FunctionCalling.SystemInstructionFile != "" {
		instruction, err := readFile(cfg.FunctionCalling.SystemInstructionFile)
		if err != nil {
			return nil, err
		}
		cfg.FunctionCalling.SystemInstruction = instruction
	}

	if cfg.Workflows.Generation.SystemInstructionFile != "" {
		instruction, err := readFile(cfg.Workflows.Generation.SystemInstructionFile)
		if err != nil {
			return nil, err
		}
		cfg.Workflows.Generation.SystemInstruction = instruction
	}

	return cfg, nil
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
