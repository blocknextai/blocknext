package getfile

import (
	"context"
	"io"
	"log/slog"
	"strings"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/drive/helpers"
)

var (
	textMimePrefixes = []string{
		"text/",
		"application/json",
		"application/xml",
		"application/javascript",
		"application/x-yaml",
		"application/yaml",
		"application/x-www-form-urlencoded",
	}

	ErrFailedToDownloadFile = apperror.Internal("failed to download file")
)

type GoogleDriveGetFileExecutorInput struct {
	FileID string `schema:"fileId"`
}

type GoogleDriveGetFileExecutor struct {
	executors.Executor
	validator   *jsonschema.Validator[GoogleDriveGetFileExecutorInput]
	fileGateway filegateway.FileGateway
}

func NewGoogleDriveGetFileExecutor(
	nodeID string,
	validator *jsonschema.Validator[GoogleDriveGetFileExecutorInput],
	fileGateway filegateway.FileGateway,
) *GoogleDriveGetFileExecutor {
	return &GoogleDriveGetFileExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator:   validator,
		fileGateway: fileGateway,
	}
}

type GoogleDriveGetFileErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

type GoogleDriveFileMetadata struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	MimeType string `json:"mimeType"`
	Size     string `json:"size"`
}

func (e *GoogleDriveGetFileExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "google_drive_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")
		client := helpers.CreateClient(ctx, accessToken)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			var metadataSuccessResponse GoogleDriveFileMetadata
			var metadataErrorResponse GoogleDriveGetFileErrorResponse
			metadataResponse, err := client.Get("/files/"+input.FileID).
				QueryParam("supportsTeamDrives", "true").
				QueryParam("supportsAllDrives", "true").
				QueryParam("fields", "id,name,mimeType,size").
				Do(&metadataSuccessResponse, &metadataErrorResponse)

			if err != nil {
				return nil, err
			}

			if !metadataResponse.IsSuccess() {
				return nil, apperror.Internal(metadataErrorResponse.Error.Message)
			}

			contentResponse, err := client.Get("/files/"+input.FileID).
				QueryParam("alt", "media").
				DoStream()

			if err != nil {
				return nil, err
			}
			defer func() {
				err := contentResponse.BodyReader.Close()
				if err != nil {
					slog.Error("Failed to close content response body reader",
						"component", "GoogleDriveGetFileExecutor",
						"error", err)
				}
			}()

			if !contentResponse.IsSuccess() {
				return nil, ErrFailedToDownloadFile
			}

			result := map[string]any{
				"metadata": map[string]any{
					"id":       metadataSuccessResponse.ID,
					"name":     metadataSuccessResponse.Name,
					"mimeType": metadataSuccessResponse.MimeType,
					"size":     metadataSuccessResponse.Size,
				},
			}

			if isTextMimeType(metadataSuccessResponse.MimeType) {
				contentBytes, err := io.ReadAll(contentResponse.BodyReader)
				if err != nil {
					return nil, err
				}
				result["file"] = string(contentBytes)
			} else {
				uploadResult, err := e.fileGateway.UploadFile(
					ctx,
					"50667bea-8f4c-4255-a54e-19d32d59927d",
					metadataSuccessResponse.Name,
					contentResponse.BodyReader,
				)
				if err != nil {
					return nil, err
				}

				result["file"] = uploadResult.URL
			}

			results = append(results, result)
		}

		return results, nil
	}
}

func isTextMimeType(mimeType string) bool {
	normalized := strings.ToLower(strings.TrimSpace(mimeType))
	for _, prefix := range textMimePrefixes {
		if strings.HasPrefix(normalized, prefix) {
			return true
		}
	}
	return false
}
