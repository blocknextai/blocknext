package createfolder

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/drive/helpers"
)

type GoogleDriveCreateFolderExecutorInput struct {
	FolderName     string `schema:"folderName"`
	ParentFolderID string `schema:"parentFolderId"`
}

type GoogleDriveCreateFolderExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[GoogleDriveCreateFolderExecutorInput]
}

func NewGoogleDriveCreateFolderExecutor(
	nodeID string,
	validator *jsonschema.Validator[GoogleDriveCreateFolderExecutorInput],
) *GoogleDriveCreateFolderExecutor {
	return &GoogleDriveCreateFolderExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

type GoogleDriveCreateFolderResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	MimeType    string   `json:"mimeType"`
	WebViewLink string   `json:"webViewLink"`
	Parents     []string `json:"parents"`
}

type GoogleDriveCreateFolderErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func (e *GoogleDriveCreateFolderExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			payload := map[string]any{
				"name":     input.FolderName,
				"mimeType": "application/vnd.google-apps.folder",
				"parents":  []string{input.ParentFolderID},
			}

			var response GoogleDriveCreateFolderResponse
			var errorResponse GoogleDriveCreateFolderErrorResponse
			res, err := client.Post("/files").
				JSONContentType().
				Body(payload).
				Do(&response, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !res.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			results = append(results, map[string]any{
				"status":      true,
				"folderId":    response.ID,
				"folderName":  response.Name,
				"webViewLink": response.WebViewLink,
			})
		}

		return results, nil
	}
}
