package deletespreadsheet

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/drive/helpers"
)

type GoogleSheetsDeleteSpreadsheetExecutorInput struct {
	SpreadsheetID string `schema:"spreadsheetId"`
}

type GoogleSheetsDeleteSpreadsheetExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[GoogleSheetsDeleteSpreadsheetExecutorInput]
}

func NewGoogleSheetsDeleteSpreadsheetExecutor(
	nodeID string,
	validator *jsonschema.Validator[GoogleSheetsDeleteSpreadsheetExecutorInput],
) *GoogleSheetsDeleteSpreadsheetExecutor {
	return &GoogleSheetsDeleteSpreadsheetExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

type GoogleSheetsDeleteSpreadsheetErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func (e *GoogleSheetsDeleteSpreadsheetExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			var errorResponse GoogleSheetsDeleteSpreadsheetErrorResponse
			response, err := client.Delete("/files/"+input.SpreadsheetID).
				Do(nil, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			results = append(results, map[string]any{
				"status": true,
			})
		}

		return results, nil
	}
}
