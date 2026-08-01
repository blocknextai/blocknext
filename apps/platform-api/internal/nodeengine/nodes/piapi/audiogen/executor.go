package audiogen

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/piapi/helpers"
)

type PiAPIAudioGenInput struct {
	Prompt               string `schema:"prompt"`
	NegativeTags         string `schema:"negativeTags"`
	GptDescriptionPrompt string `schema:"gptDescriptionPrompt"`
	Title                string `schema:"title"`
	LyricsType           string `schema:"lyricsType"`
	Seed                 int    `schema:"seed"`
	Lyrics               string `schema:"lyrics"`
}

type PiAPIAudioGenExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[PiAPIAudioGenInput]
}

func NewPiAPIAudioGenExecutor(
	nodeID string,
	validator *jsonschema.Validator[PiAPIAudioGenInput],
) *PiAPIAudioGenExecutor {
	return &PiAPIAudioGenExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

type PiAPIAudioGenSuccessResponse struct {
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

type PiAPIAudioGenErrorResponse struct {
	Message string `json:"message"`
}

func (e *PiAPIAudioGenExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			payload := generatePayload(*input)

			var successResponse PiAPIAudioGenSuccessResponse
			var errorResponse PiAPIAudioGenErrorResponse
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

			audioData, err := WaitAudioGen(ctx, client, WaitAudioGenInput{
				TaskID: successResponse.Data.TaskID,
			})
			if err != nil {
				return nil, err
			}

			results = append(results, audioData...)
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
	Prompt               string `json:"prompt,omitempty"`
	NegativeTags         string `json:"negative_tags,omitempty"`
	GptDescriptionPrompt string `json:"gpt_description_prompt,omitempty"`
	Title                string `json:"title,omitempty"`
	LyricsType           string `json:"lyrics_type,omitempty"`
	Seed                 int    `json:"seed,omitempty"`
	Lyrics               string `json:"lyrics,omitempty"`
}

func generatePayload(input PiAPIAudioGenInput) *Payload {
	return &Payload{
		Model:    "music-u",
		TaskType: "generate_music",
		Input: &InputPayload{
			Prompt:               input.Prompt,
			NegativeTags:         input.NegativeTags,
			GptDescriptionPrompt: input.GptDescriptionPrompt,
			Title:                input.Title,
			LyricsType:           input.LyricsType,
			Seed:                 input.Seed,
			Lyrics:               input.Lyrics,
		},
	}
}
