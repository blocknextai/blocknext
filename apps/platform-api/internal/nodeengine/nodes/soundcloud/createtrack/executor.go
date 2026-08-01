package createtrack

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/soundcloud/helpers"
)

var (
	ErrExecutorEmptyResponse = apperror.Internal("empty response")
	ErrUploadFailed          = apperror.Internal("upload failed")
)

type SoundCloudCreateTrackExecutorInput struct {
	MP3Link     string `schema:"mp3Link"`
	Title       string `schema:"title"`
	Description string `schema:"description"`
}

type SoundCloudCreateTrackExecutor struct {
	executors.Executor
	validator   *jsonschema.Validator[SoundCloudCreateTrackExecutorInput]
	fileGateway filegateway.FileGateway
}

func NewSoundCloudCreateTrackExecutor(
	nodeID string,
	validator *jsonschema.Validator[SoundCloudCreateTrackExecutorInput],
	fileGateway filegateway.FileGateway,
) *SoundCloudCreateTrackExecutor {
	return &SoundCloudCreateTrackExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator:   validator,
		fileGateway: fileGateway,
	}
}

func (e *SoundCloudCreateTrackExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "soundcloud_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")
		client := helpers.CreateClient(ctx, accessToken)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			downloadResult, err := e.fileGateway.DownloadFile(ctx, input.MP3Link)
			if err != nil {
				return nil, err
			}
			defer func() {
				_ = downloadResult.DataReader.Close()
			}()

			uploadResponse, err := e.uploadToSoundCloud(client, *input, downloadResult)
			if err != nil {
				return nil, err
			}

			if uploadResponse == nil {
				return nil, ErrExecutorEmptyResponse
			}

			results = append(results, map[string]any{
				"id": uploadResponse.ID,
			})
		}

		return results, nil
	}
}

type UploadResponse struct {
	ID           int    `json:"id"`
	PermalinkURL string `json:"permalink_url"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Sharing      string `json:"sharing"`
	State        string `json:"state"`
}

type ErrorResponse struct {
	Errors  []map[string]any `json:"errors"`
	Message string           `json:"message"`
}

func (e *SoundCloudCreateTrackExecutor) uploadToSoundCloud(client *httpclient.Client, input SoundCloudCreateTrackExecutorInput, downloadResult *filegateway.DownloadResult) (*UploadResponse, error) {
	filename := downloadResult.Filename

	multipartBuilder := client.Post("/tracks").MultipartFormBody().
		AddFileReader("track[asset_data]", filename, downloadResult.DataReader).
		AddField("track[title]", input.Title)

	if input.Description != "" {
		multipartBuilder.AddField("track[description]", input.Description)
	}

	var successResponse UploadResponse
	var errorResponse ErrorResponse
	response, err := multipartBuilder.Do(&successResponse, &errorResponse)

	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		return nil, ErrUploadFailed
	}

	return &successResponse, nil
}
