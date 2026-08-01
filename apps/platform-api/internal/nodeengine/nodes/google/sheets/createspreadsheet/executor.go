package createspreadsheet

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/drive/helpers"
)

var (
	ErrExecutorEmptyResponse = apperror.Internal("empty response")
)

type GoogleSheetsCreateSpreadsheetExecutorInput struct {
	Title string `schema:"title"`
}

type GoogleSheetsCreateSpreadsheetExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[GoogleSheetsCreateSpreadsheetExecutorInput]
}

func NewGoogleSheetsCreateSpreadsheetExecutor(
	nodeID string,
	validator *jsonschema.Validator[GoogleSheetsCreateSpreadsheetExecutorInput],
) *GoogleSheetsCreateSpreadsheetExecutor {
	return &GoogleSheetsCreateSpreadsheetExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

type GoogleSheetsCreateSpreadsheetSuccessResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	WebViewLink string `json:"webViewLink"`
}

type GoogleSheetsCreateSpreadsheetErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func (e *GoogleSheetsCreateSpreadsheetExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "google_sheets_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")
		client := helpers.CreateClient(ctx, accessToken)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			var successResponse GoogleSheetsCreateSpreadsheetSuccessResponse
			var errorResponse GoogleSheetsCreateSpreadsheetErrorResponse
			response, err := client.Post("/files").
				JSONContentType().
				Body(map[string]any{
					"name":     input.Title,
					"mimeType": "application/vnd.google-apps.spreadsheet",
				}).
				Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			if successResponse.ID == "" {
				return nil, ErrExecutorEmptyResponse
			}

			results = append(results, map[string]any{
				"status":      true,
				"id":          successResponse.ID,
				"name":        successResponse.Name,
				"webViewLink": successResponse.WebViewLink,
			})
		}

		return results, nil
	}
}
