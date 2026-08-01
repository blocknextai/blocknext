package videogen

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/piapi/helpers"
)

type PiAPIVideoGenInput struct {
	Prompt                         string  `schema:"prompt"`
	NegativePrompt                 string  `schema:"negativePrompt"`
	CfgScale                       float64 `schema:"cfgScale"`
	Duration                       int     `schema:"duration"`
	AspectRatio                    string  `schema:"aspectRatio"`
	CameraControlType              string  `schema:"cameraControlType"`
	CameraControlConfigHorizontal  int     `schema:"cameraControlConfigHorizontal"`
	CameraControlConfigVertical    int     `schema:"cameraControlConfigVertical"`
	CameraControlConfigPan         int     `schema:"cameraControlConfigPan"`
	CameraControlConfigTilt        int     `schema:"cameraControlConfigTilt"`
	CameraControlConfigRoll        int     `schema:"cameraControlConfigRoll"`
	CameraControlConfigZoom        int     `schema:"cameraControlConfigZoom"`
	Mode                           string  `schema:"mode"`
	Version                        string  `schema:"version"`
	ImageURL                       string  `schema:"imageUrl"`
	ImageTailURL                   string  `schema:"imageTailUrl"`
	MotionBrushMaskURL             string  `schema:"motionBrushMaskUrl"`
	MotionBrushPointsStaticMasksX  int     `schema:"motionBrushPointsStaticMasksX"`
	MotionBrushPointsStaticMasksY  int     `schema:"motionBrushPointsStaticMasksY"`
	MotionBrushPointsDynamicMasksX int     `schema:"motionBrushPointsDynamicMasksX"`
	MotionBrushPointsDynamicMasksY int     `schema:"motionBrushPointsDynamicMasksY"`
	Effect                         string  `schema:"effect"`
	ServiceMode                    string  `schema:"serviceMode"`
}

type PiAPIVideoGenExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[PiAPIVideoGenInput]
}

func NewPiAPIVideoGenExecutor(
	nodeID string,
	validator *jsonschema.Validator[PiAPIVideoGenInput],
) *PiAPIVideoGenExecutor {
	return &PiAPIVideoGenExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

type PiAPIVideoGenSuccessResponse struct {
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
	} `json:"data"`
	Message string `json:"message"`
}

type PiAPIVideoGenErrorResponse struct {
	Message string `json:"message"`
}

func (e *PiAPIVideoGenExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			var successResponse PiAPIVideoGenSuccessResponse
			var errorResponse PiAPIVideoGenErrorResponse
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

			videoData, err := WaitVideoGen(ctx, client, WaitVideoGenInput{
				TaskID: successResponse.Data.TaskID,
			})
			if err != nil {
				return nil, err
			}

			results = append(results, videoData...)
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
	Prompt         string         `json:"prompt,omitempty"`
	NegativePrompt string         `json:"negative_prompt,omitempty"`
	CfgScale       float64        `json:"cfg_scale,omitempty"`
	Duration       int            `json:"duration,omitempty"`
	AspectRatio    string         `json:"aspect_ratio,omitempty"`
	CameraControl  *CameraControl `json:"camera_control,omitempty"`
	Mode           string         `json:"mode,omitempty"`
	Version        string         `json:"version,omitempty"`
	ImageURL       string         `json:"image_url,omitempty"`
	ImageTailURL   string         `json:"image_tail_url,omitempty"`
	MotionBrush    *MotionBrush   `json:"motion_brush,omitempty"`
	Effect         string         `json:"effect,omitempty"`
	Config         *Config        `json:"config,omitempty"`
}

type CameraControl struct {
	Type   string        `json:"type,omitempty"`
	Config *CameraConfig `json:"config,omitempty"`
}

type CameraConfig struct {
	Horizontal int `json:"horizontal,omitempty"`
	Vertical   int `json:"vertical,omitempty"`
	Pan        int `json:"pan,omitempty"`
	Tilt       int `json:"tilt,omitempty"`
	Roll       int `json:"roll,omitempty"`
	Zoom       int `json:"zoom,omitempty"`
}

type MotionBrush struct {
	MaskURL      string      `json:"mask_url,omitempty"`
	StaticMasks  *MaskPoints `json:"static_masks,omitempty"`
	DynamicMasks *MaskPoints `json:"dynamic_masks,omitempty"`
}

type MaskPoints struct {
	Points *MotionPoints `json:"points,omitempty"`
}

type MotionPoints struct {
	X int `json:"x,omitempty"`
	Y int `json:"y,omitempty"`
}

type Config struct {
	ServiceMode string `json:"service_mode,omitempty"`
}

func generatePayload(input *PiAPIVideoGenInput) *Payload {
	payload := &Payload{
		Model:    "kling",
		TaskType: "video_generation",
		Input: &InputPayload{
			Prompt:         input.Prompt,
			NegativePrompt: input.NegativePrompt,
			CfgScale:       input.CfgScale,
			Duration:       input.Duration,
			AspectRatio:    input.AspectRatio,
			Mode:           input.Mode,
			Version:        input.Version,
			ImageURL:       input.ImageURL,
			ImageTailURL:   input.ImageTailURL,
			Effect:         input.Effect,
		},
	}

	if input.CameraControlType != "" ||
		input.CameraControlConfigHorizontal != 0 || input.CameraControlConfigVertical != 0 ||
		input.CameraControlConfigPan != 0 || input.CameraControlConfigTilt != 0 ||
		input.CameraControlConfigRoll != 0 || input.CameraControlConfigZoom != 0 {

		cameraConfig := &CameraConfig{}
		if input.CameraControlConfigHorizontal != 0 {
			cameraConfig.Horizontal = input.CameraControlConfigHorizontal
		}
		if input.CameraControlConfigVertical != 0 {
			cameraConfig.Vertical = input.CameraControlConfigVertical
		}
		if input.CameraControlConfigPan != 0 {
			cameraConfig.Pan = input.CameraControlConfigPan
		}
		if input.CameraControlConfigTilt != 0 {
			cameraConfig.Tilt = input.CameraControlConfigTilt
		}
		if input.CameraControlConfigRoll != 0 {
			cameraConfig.Roll = input.CameraControlConfigRoll
		}
		if input.CameraControlConfigZoom != 0 {
			cameraConfig.Zoom = input.CameraControlConfigZoom
		}

		payload.Input.CameraControl = &CameraControl{
			Type:   input.CameraControlType,
			Config: cameraConfig,
		}
	}

	if input.MotionBrushMaskURL != "" ||
		input.MotionBrushPointsStaticMasksX != 0 || input.MotionBrushPointsStaticMasksY != 0 ||
		input.MotionBrushPointsDynamicMasksX != 0 || input.MotionBrushPointsDynamicMasksY != 0 {

		motionBrush := &MotionBrush{
			MaskURL: input.MotionBrushMaskURL,
		}

		if input.MotionBrushPointsStaticMasksX != 0 || input.MotionBrushPointsStaticMasksY != 0 {
			motionBrush.StaticMasks = &MaskPoints{
				Points: &MotionPoints{
					X: input.MotionBrushPointsStaticMasksX,
					Y: input.MotionBrushPointsStaticMasksY,
				},
			}
		}

		if input.MotionBrushPointsDynamicMasksX != 0 || input.MotionBrushPointsDynamicMasksY != 0 {
			motionBrush.DynamicMasks = &MaskPoints{
				Points: &MotionPoints{
					X: input.MotionBrushPointsDynamicMasksX,
					Y: input.MotionBrushPointsDynamicMasksY,
				},
			}
		}

		payload.Input.MotionBrush = motionBrush
	}

	if input.ServiceMode != "" {
		payload.Input.Config = &Config{
			ServiceMode: input.ServiceMode,
		}
	}

	return payload
}
