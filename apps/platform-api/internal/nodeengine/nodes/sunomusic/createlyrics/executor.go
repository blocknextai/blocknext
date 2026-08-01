package createlyrics

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/sunomusic/helpers"
)

type SunoMusicCreateLyricsInput struct {
	Prompt string `schema:"prompt"`
}

type SunoMusicCreateLyricsExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[SunoMusicCreateLyricsInput]
}

func NewSunoMusicCreateLyricsExecutor(
	nodeID string,
	validator *jsonschema.Validator[SunoMusicCreateLyricsInput],
) *SunoMusicCreateLyricsExecutor {
	return &SunoMusicCreateLyricsExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

type SunoMusicCreateLyricsSuccessResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		TaskID string `json:"taskId"`
	} `json:"data"`
}

type SunoMusicCreateLyricsErrorResponse struct{}

func (e *SunoMusicCreateLyricsExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			var successResponse SunoMusicCreateLyricsSuccessResponse
			var errorResponse SunoMusicCreateLyricsErrorResponse
			_, err = client.Post("/lyrics").
				JSONContentType().
				Body(map[string]any{
					"prompt":      input.Prompt,
					"callbackUrl": "https://api.example.com/callback",
				}).
				Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if successResponse.Code != 200 {
				return nil, apperror.Internal(successResponse.Msg)
			}

			lyricsData, err := WaitLyrics(ctx, client, WaitLyricsInput{
				TaskID: successResponse.Data.TaskID,
			})
			if err != nil {
				return nil, err
			}

			results = append(results, lyricsData)
		}

		return results, nil
	}
}
