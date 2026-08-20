package createtextfile

import (
	"bytes"
	"context"
	"strconv"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
)

type GoogleDriveCreateTextFileExecutorInput struct {
	FileName       string `schema:"fileName"`
	Content        string `schema:"content"`
	ParentFolderID string `schema:"parentFolderId"`
}

type GoogleDriveCreateTextFileExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[GoogleDriveCreateTextFileExecutorInput]
}

func NewGoogleDriveCreateTextFileExecutor(
	nodeID string,
	validator *jsonschema.Validator[GoogleDriveCreateTextFileExecutorInput],
) *GoogleDriveCreateTextFileExecutor {
	return &GoogleDriveCreateTextFileExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

type GoogleDriveCreateTextFileResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	MimeType    string   `json:"mimeType"`
	WebViewLink string   `json:"webViewLink"`
	Parents     []string `json:"parents"`
}

type GoogleDriveCreateTextFileErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func (e *GoogleDriveCreateTextFileExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			payload := map[string]any{
				"name":     input.FileName,
				"mimeType": "text/plain",
				"parents":  []string{input.ParentFolderID},
			}

			client := httpclient.NewClientBuilder().
				Context(ctx).
				BearerAuth(accessToken).
				BaseURL("https://www.googleapis.com/upload/drive/v3").
				Build()

			var sessionResponse GoogleDriveCreateTextFileResponse
			var sessionError GoogleDriveCreateTextFileErrorResponse

			initReq := client.Post("/files?uploadType=resumable").
				JSONContentType().
				Header("X-Upload-Content-Type", "text/plain").
				Body(payload)

			initRes, err := initReq.Do(&sessionResponse, &sessionError)

			if err != nil {
				return nil, err
			}

			if !initRes.IsSuccess() {
				return nil, apperror.Internal(sessionError.Error.Message)
			}

			uploadURL := initRes.Headers.Get("Location")
			if uploadURL == "" {
				return nil, apperror.Internal(sessionError.Error.Message)
			}

			contentBytes := []byte(input.Content)
			putClient := httpclient.NewClientBuilder().
				Context(ctx).
				BearerAuth(accessToken).
				BaseURL(uploadURL).
				Build()

			var finalResponse GoogleDriveCreateTextFileResponse
			var finalError GoogleDriveCreateTextFileErrorResponse

			res, err := putClient.Put("").
				Header("Content-Length", strconv.Itoa(len(contentBytes))).
				Header("Content-Type", "text/plain").
				BodyReader(bytes.NewReader(contentBytes)).
				Do(&finalResponse, &finalError)

			if err != nil {
				return nil, err
			}

			if !res.IsSuccess() {
				return nil, apperror.Internal(finalError.Error.Message)
			}

			results = append(results, map[string]any{
				"status":      true,
				"fileId":      finalResponse.ID,
				"fileName":    finalResponse.Name,
				"webViewLink": finalResponse.WebViewLink,
			})
		}

		return results, nil
	}
}
