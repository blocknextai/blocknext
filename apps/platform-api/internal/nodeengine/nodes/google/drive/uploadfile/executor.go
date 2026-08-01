package uploadfile

import (
	"bytes"
	"context"
	"io"
	"strconv"
	"strings"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
)

const (
	uploadChunkSize = 2 * 1024 * 1024
)

var (
	errFailedToDownloadFile         = apperror.Internal("failed to download file")
	errFailedToReadFileContent      = apperror.Internal("failed to read file content")
	errFailedToStartResumableUpload = apperror.Internal("failed to start resumable upload session")
	errMissingUploadURL             = apperror.Internal("missing upload url in response")
	errChunkUploadFailed            = apperror.Internal("chunk upload failed")
)

type GoogleDriveUploadFileExecutorInput struct {
	FolderID string   `schema:"folderId"`
	Files    []string `schema:"files"`
}

type GoogleDriveUploadFileExecutor struct {
	executors.Executor
	validator   *jsonschema.Validator[GoogleDriveUploadFileExecutorInput]
	fileGateway filegateway.FileGateway
}

func NewGoogleDriveUploadFileExecutor(
	nodeID string,
	validator *jsonschema.Validator[GoogleDriveUploadFileExecutorInput],
	fileGateway filegateway.FileGateway,
) *GoogleDriveUploadFileExecutor {
	return &GoogleDriveUploadFileExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator:   validator,
		fileGateway: fileGateway,
	}
}

func (e *GoogleDriveUploadFileExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "google_drive_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			uploadedFiles := make([]map[string]any, 0)

			for _, fileURL := range input.Files {
				downloadResult, err := e.fileGateway.DownloadFile(ctx, fileURL)
				if err != nil {
					return nil, errFailedToDownloadFile
				}
				defer func() {
					_ = downloadResult.DataReader.Close()
				}()

				fileData, err := io.ReadAll(downloadResult.DataReader)
				if err != nil {
					return nil, errFailedToReadFileContent
				}

				contentType := downloadResult.ContentType
				if contentType == "" {
					contentType = "application/octet-stream"
				}

				fileName := downloadResult.Filename

				metadata := map[string]any{
					"name":     fileName,
					"mimeType": contentType,
					"parents":  []string{input.FolderID},
				}

				client := httpclient.NewClientBuilder().
					Context(ctx).
					BearerAuth(accessToken).
					BaseURL("https://www.googleapis.com/upload/drive/v3").
					Build()

				var sessionResponse map[string]any
				var sessionError map[string]any

				initReq := client.Post("/files?uploadType=resumable").
					JSONContentType().
					Header("X-Upload-Content-Type", contentType).
					Body(metadata)

				initRes, err := initReq.Do(&sessionResponse, &sessionError)

				if err != nil {
					return nil, errFailedToStartResumableUpload
				}

				if initRes.IsError() {
					return nil, errFailedToStartResumableUpload
				}

				uploadURL := initRes.Headers.Get("Location")
				if uploadURL == "" {
					return nil, errMissingUploadURL
				}

				totalSize := len(fileData)
				putClient := httpclient.NewClientBuilder().
					Context(ctx).
					BearerAuth(accessToken).
					BaseURL(uploadURL).
					Build()

				var finalResponse map[string]any
				var finalError map[string]any

				for start := 0; start < totalSize; start += uploadChunkSize {
					end := min(start+uploadChunkSize, totalSize)

					chunk := fileData[start:end]

					var rb strings.Builder
					rb.WriteString("bytes ")
					rb.WriteString(strconv.Itoa(start))
					rb.WriteByte('-')
					rb.WriteString(strconv.Itoa(end - 1))
					rb.WriteByte('/')
					rb.WriteString(strconv.Itoa(totalSize))
					rangeHeader := rb.String()

					res, err := putClient.Put("").
						Header("Content-Length", strconv.Itoa(len(chunk))).
						Header("Content-Type", contentType).
						Header("Content-Range", rangeHeader).
						BodyReader(bytes.NewReader(chunk)).
						Do(&finalResponse, &finalError)

					if err != nil {
						return nil, errChunkUploadFailed
					}

					if end == totalSize && res.IsSuccess() && finalResponse != nil {
						uploadedFiles = append(uploadedFiles, map[string]any{
							"id":          finalResponse["id"],
							"name":        finalResponse["name"],
							"mimeType":    finalResponse["mimeType"],
							"webViewLink": finalResponse["webViewLink"],
						})
					}
				}
			}

			results = append(results, map[string]any{
				"status": len(uploadedFiles) == len(input.Files),
			})
		}

		return results, nil
	}
}
