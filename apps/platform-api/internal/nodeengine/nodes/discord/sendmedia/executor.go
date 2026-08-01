package sendmedia

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/discord/helpers"
)

var (
	ErrFileTooLarge = apperror.Internal("file too large")
)

const (
	MaxFileSize = 8 * 1024 * 1024
)

type DiscordSendMediaExecutorInput struct {
	ChannelID string   `schema:"channelId"`
	MediaURLs []string `schema:"mediaUrls"`
}

type DiscordSendMediaExecutor struct {
	executors.Executor
	validator   *jsonschema.Validator[DiscordSendMediaExecutorInput]
	fileGateway filegateway.FileGateway
}

func NewDiscordSendMediaExecutor(
	nodeID string,
	validator *jsonschema.Validator[DiscordSendMediaExecutorInput],
	fileGateway filegateway.FileGateway,
) *DiscordSendMediaExecutor {
	return &DiscordSendMediaExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator:   validator,
		fileGateway: fileGateway,
	}
}

type SuccessResponse struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	Timestamp string `json:"timestamp"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *DiscordSendMediaExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "discord_api")
		botToken := credential.String("botToken")
		client := helpers.CreateClient(ctx, botToken)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			for _, mediaURL := range input.MediaURLs {
				messageID, err := e.sendSingleMedia(ctx, client, input.ChannelID, mediaURL)
				if err != nil {
					results = append(results, map[string]any{
						"status": false,
						"error":  err.Error(),
					})
					continue
				}
				results = append(results, map[string]any{
					"status":    true,
					"messageId": messageID,
				})
			}
		}

		return results, nil
	}
}

func (e *DiscordSendMediaExecutor) sendSingleMedia(ctx context.Context, client *httpclient.Client, channelID, mediaURL string) (string, error) {
	downloadResult, err := e.fileGateway.DownloadFile(ctx, mediaURL)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = downloadResult.DataReader.Close()
	}()

	if downloadResult.Size > MaxFileSize {
		return "", ErrFileTooLarge
	}

	filename := downloadResult.Filename

	var successResponse SuccessResponse
	var errorResponse ErrorResponse
	response, err := client.Post("/channels/"+channelID+"/messages").
		MultipartFormBody().
		AddFileReader("files[0]", filename, downloadResult.DataReader).
		Do(&successResponse, &errorResponse)

	if err != nil {
		return "", err
	}

	if !response.IsSuccess() {
		return "", apperror.Internal(errorResponse.Message)
	}

	return successResponse.ID, nil
}
