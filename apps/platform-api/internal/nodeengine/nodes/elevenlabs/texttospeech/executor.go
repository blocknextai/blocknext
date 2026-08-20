package texttospeech

import (
	"bytes"
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/elevenlabs/helpers"
)

var (
	ErrGenerateElevenlabsTextToSpeech = apperror.Internal("generate error")
)

type ElevenlabsTextToSpeechExecutorInput struct {
	Text         string `schema:"text"`
	VoiceID      string `schema:"voiceId"`
	ModelID      string `schema:"modelId"`
	OutputFormat string `schema:"outputFormat"`
}

type ElevenlabsTextToSpeechExecutor struct {
	executors.Executor
	validator   *jsonschema.Validator[ElevenlabsTextToSpeechExecutorInput]
	fileGateway filegateway.FileGateway
}

func NewElevenlabsTextToSpeechExecutor(
	nodeID string,
	validator *jsonschema.Validator[ElevenlabsTextToSpeechExecutorInput],
	fileGateway filegateway.FileGateway,
) *ElevenlabsTextToSpeechExecutor {
	return &ElevenlabsTextToSpeechExecutor{
		ID:          nodeID,
		validator:   validator,
		fileGateway: fileGateway,
	}
}

func (e *ElevenlabsTextToSpeechExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "elevenlabs_api")
		apiKey := credential.String("apiKey")
		client := helpers.CreateClient(ctx, apiKey)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			response, err := client.Post("/text-to-speech/"+input.VoiceID).
				QueryParam("output_format", input.OutputFormat).
				JSONContentType().
				Body(map[string]any{
					"text":     input.Text,
					"model_id": input.ModelID,
				}).
				Do(nil, nil)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, ErrGenerateElevenlabsTextToSpeech
			}

			result, err := e.fileGateway.UploadFile(
				ctx,
				"8044ae7e-4d97-4531-b945-94f824d71987",
				"elevenlabs-texttospeech.mp3",
				bytes.NewReader(response.Body),
			)

			if err != nil {
				return nil, err
			}

			results = append(results, map[string]any{
				"audio": result.URL,
			})
		}

		return results, nil
	}
}
