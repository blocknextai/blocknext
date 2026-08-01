package sendmedia

import (
	"context"
	"io"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/cast"
	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/slack/helpers"
)

var (
	ErrFileTooLarge     = apperror.Internal("file too large")
	ErrFileUploadFailed = apperror.Internal("file upload failed")
)

const (
	MaxFileSize = 1000 * 1024 * 1024
)

type SlackSendMediaExecutorInput struct {
	Channel   string   `schema:"channel"`
	MediaURLs []string `schema:"mediaUrls"`
	Text      string   `schema:"text"`
}

type SlackSendMediaExecutor struct {
	executors.Executor
	validator   *jsonschema.Validator[SlackSendMediaExecutorInput]
	fileGateway filegateway.FileGateway
}

func NewSlackSendMediaExecutor(
	nodeID string,
	validator *jsonschema.Validator[SlackSendMediaExecutorInput],
	fileGateway filegateway.FileGateway,
) *SlackSendMediaExecutor {
	return &SlackSendMediaExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator:   validator,
		fileGateway: fileGateway,
	}
}

func (e *SlackSendMediaExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "slack_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")
		client := helpers.CreateClient(ctx, accessToken)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			for _, mediaURL := range input.MediaURLs {
				if _, err := e.sendSingleMedia(ctx, client, input.Channel, mediaURL, input.Text); err != nil {
					results = append(results, map[string]any{
						"status": false,
						"error":  err.Error(),
					})
					continue
				}
				results = append(results, map[string]any{
					"status": true,
				})
			}
		}

		return results, nil
	}
}

type GetUploadURLResponse struct {
	Ok        bool   `json:"ok"`
	UploadURL string `json:"upload_url"`
	FileID    string `json:"file_id"`
	Error     string `json:"error"`
}

type CompleteUploadResponse struct {
	Ok    bool `json:"ok"`
	Files []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"files"`
	Error string `json:"error"`
}

func (e *SlackSendMediaExecutor) sendSingleMedia(ctx context.Context, client *httpclient.Client, channel, mediaURL, text string) (string, error) {
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

	uploadURL, fileID, err := e.getUploadURL(client, filename, downloadResult.Size)
	if err != nil {
		return "", err
	}

	err = e.uploadFileToURL(ctx, uploadURL, filename, downloadResult.DataReader)
	if err != nil {
		return "", err
	}

	messageID, err := e.completeUpload(client, fileID, filename, channel, text)
	if err != nil {
		return "", err
	}

	return messageID, nil
}

func (e *SlackSendMediaExecutor) getUploadURL(client *httpclient.Client, filename string, fileSize int64) (string, string, error) {
	var response GetUploadURLResponse
	_, err := client.Get("/files.getUploadURLExternal").
		QueryParam("filename", filename).
		QueryParam("length", cast.ToString(fileSize)).
		Do(&response, nil)

	if err != nil {
		return "", "", err
	}

	if !response.Ok {
		return "", "", apperror.Internal(response.Error)
	}

	return response.UploadURL, response.FileID, nil
}

func (e *SlackSendMediaExecutor) uploadFileToURL(ctx context.Context, uploadURL, filename string, fileContent io.Reader) error {
	uploadClient := httpclient.NewClientBuilder().
		Context(ctx).
		Build()

	response, err := uploadClient.Post(uploadURL).
		MultipartFormBody().
		AddFileReader("file", filename, fileContent).
		Do(nil, nil)

	if err != nil {
		return err
	}

	if !response.IsSuccess() {
		return ErrFileUploadFailed
	}

	return nil
}

func (e *SlackSendMediaExecutor) completeUpload(client *httpclient.Client, fileID, filename, channel, text string) (string, error) {
	payload := map[string]any{
		"channel_id": channel,
		"files": []map[string]string{
			{
				"id":    fileID,
				"title": filename,
			},
		},
	}

	if text != "" {
		payload["initial_comment"] = text
	}

	var response CompleteUploadResponse
	_, err := client.Post("/files.completeUploadExternal").
		JSONContentType().
		Body(payload).
		Do(&response, nil)

	if err != nil {
		return "", err
	}

	if !response.Ok {
		return "", apperror.Internal(response.Error)
	}

	if len(response.Files) == 0 {
		return "", nil
	}

	return response.Files[0].ID, nil
}
