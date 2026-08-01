package createmusic

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/sunomusic/helpers"
)

type SunoMusicCreateMusicInput struct {
	Prompt       string `schema:"prompt"`
	Style        string `schema:"style"`
	Title        string `schema:"title"`
	CustomMode   bool   `schema:"customMode"`
	Instrumental bool   `schema:"instrumental"`
	NegativeTags string `schema:"negativeTags"`
	Model        string `schema:"model"`
}

type SunoMusicCreateMusicExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[SunoMusicCreateMusicInput]
}

func NewSunoMusicCreateMusicExecutor(
	nodeID string,
	validator *jsonschema.Validator[SunoMusicCreateMusicInput],
) *SunoMusicCreateMusicExecutor {
	return &SunoMusicCreateMusicExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

type SunoMusicCreateMusicSuccessResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		TaskID string `json:"taskId"`
	} `json:"data"`
}

type SunoMusicCreateMusicErrorResponse struct{}

func (e *SunoMusicCreateMusicExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "sunomusic_api")
		apiKey := credential.String("apiKey")
		client := helpers.CreateClient(ctx, apiKey)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			var successResponse SunoMusicCreateMusicSuccessResponse
			var errorResponse SunoMusicCreateMusicErrorResponse
			_, err = client.Post("/generate").
				JSONContentType().
				Body(map[string]any{
					"prompt":       input.Prompt,
					"style":        input.Style,
					"title":        input.Title,
					"customMode":   input.CustomMode,
					"instrumental": input.Instrumental,
					"negativeTags": input.NegativeTags,
					"model":        input.Model,
					"callbackUrl":  "https://api.example.com/callback",
				}).
				Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if successResponse.Code != 200 {
				return nil, apperror.Internal(successResponse.Msg)
			}

			musicData, err := WaitMusic(ctx, client, WaitMusicInput{
				TaskID: successResponse.Data.TaskID,
			})
			if err != nil {
				return nil, err
			}

			results = append(results, musicData...)
		}

		return results, nil
	}
}
