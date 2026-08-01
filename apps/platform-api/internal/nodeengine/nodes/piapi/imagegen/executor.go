package imagegen

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/piapi/helpers"
)

type PiAPIImageGenInput struct {
	Model           string  `schema:"model"`
	Prompt          string  `schema:"prompt"`
	Denoise         float64 `schema:"denoise"`
	NegativePrompt  string  `schema:"negativePrompt"`
	GuidanceScale   float64 `schema:"guidanceScale"`
	Width           float64 `schema:"width"`
	Height          float64 `schema:"height"`
	BatchSize       float64 `schema:"batchSize"`
	LoraType        string  `schema:"loraType"`
	LoraImage       string  `schema:"loraImage"`
	ControlNetType  string  `schema:"controlNetType"`
	ControlNetImage string  `schema:"controlNetImage"`
	ServiceMode     string  `schema:"serviceMode"`
}

type PiAPIImageGenExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[PiAPIImageGenInput]
}

func NewPiAPIImageGenExecutor(
	nodeID string,
	validator *jsonschema.Validator[PiAPIImageGenInput],
) *PiAPIImageGenExecutor {
	return &PiAPIImageGenExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

type PiAPIImageGenSuccessResponse struct {
	Code int `json:"code"`
	Data struct {
		TaskID   string `json:"task_id"`
		Model    string `json:"model"`
		TaskType string `json:"task_type"`
		Status   string `json:"status"`
		Input    any    `json:"input"`
		Output   any    `json:"output"`
		Meta     struct {
			CreatedAt string `json:"created_at"`
			StartedAt string `json:"started_at"`
			EndedAt   string `json:"ended_at"`
			Usage     struct {
				Type    string  `json:"type"`
				Frozen  float64 `json:"frozen"`
				Consume float64 `json:"consume"`
			} `json:"usage"`
			IsUsingPrivatePool bool `json:"is_using_private_pool"`
		} `json:"meta"`
		Detail any   `json:"detail"`
		Logs   []any `json:"logs"`
		Error  struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	} `json:"data"`
	Message string `json:"message"`
}

type PiAPIImageGenErrorResponse struct {
	Message string `json:"message"`
}

func (e *PiAPIImageGenExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "piapi_api")
		apiKey := credential.String("apiKey")
		client := helpers.CreateClient(ctx, apiKey)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			payload := generatePayload(input)

			var successResponse PiAPIImageGenSuccessResponse
			var errorResponse PiAPIImageGenErrorResponse
			response, err := client.Post("/task").
				JSONContentType().
				Body(payload).
				Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Message)
			}

			imageData, err := WaitImageGen(ctx, client, WaitImageGenInput{
				TaskID: successResponse.Data.TaskID,
			})
			if err != nil {
				return nil, err
			}

			results = append(results, imageData)
		}

		return results, nil
	}
}

type Payload struct {
	Model    string        `json:"model,omitempty"`
	TaskType string        `json:"task_type,omitempty"`
	Input    *InputPayload `json:"input,omitempty"`
}

type InputPayload struct {
	Prompt             string              `json:"prompt,omitempty"`
	Denoise            float64             `json:"denoise,omitempty"`
	NegativePrompt     string              `json:"negative_prompt,omitempty"`
	GuidanceScale      float64             `json:"guidance_scale,omitempty"`
	Width              float64             `json:"width,omitempty"`
	Height             float64             `json:"height,omitempty"`
	BatchSize          float64             `json:"batch_size,omitempty"`
	LoraSettings       *LoraSettings       `json:"lora_settings,omitempty"`
	ControlNetSettings *ControlNetSettings `json:"control_net_settings,omitempty"`
	Config             *Config             `json:"config,omitempty"`
}

type LoraSettings struct {
	LoraType  string `json:"lora_type,omitempty"`
	LoraImage string `json:"lora_image,omitempty"`
}

type ControlNetSettings struct {
	ControlNetType  string `json:"control_type,omitempty"`
	ControlNetImage string `json:"control_image,omitempty"`
}

type Config struct {
	ServiceMode string `json:"service_mode,omitempty"`
}

func generatePayload(input *PiAPIImageGenInput) *Payload {
	payload := &Payload{
		Model:    input.Model,
		TaskType: "txt2img",
		Input: &InputPayload{
			Prompt:         input.Prompt,
			Denoise:        input.Denoise,
			NegativePrompt: input.NegativePrompt,
			GuidanceScale:  input.GuidanceScale,
			Width:          input.Width,
			Height:         input.Height,
			BatchSize:      input.BatchSize,
		},
	}

	if input.LoraType != "" || input.LoraImage != "" {
		payload.Input.LoraSettings = &LoraSettings{
			LoraType:  input.LoraType,
			LoraImage: input.LoraImage,
		}
	}

	if input.ControlNetType != "" || input.ControlNetImage != "" {
		payload.Input.ControlNetSettings = &ControlNetSettings{
			ControlNetType:  input.ControlNetType,
			ControlNetImage: input.ControlNetImage,
		}
	}

	if input.ServiceMode != "" {
		payload.Input.Config = &Config{
			ServiceMode: input.ServiceMode,
		}
	}

	return payload
}
