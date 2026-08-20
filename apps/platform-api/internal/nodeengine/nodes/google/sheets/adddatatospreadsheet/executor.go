package adddatatospreadsheet

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/sheets/helpers"
)

type GoogleSheetsAddDataToSpreadsheetExecutorInput struct {
	SpreadsheetID    string     `schema:"spreadsheetId"`
	Range            string     `schema:"range"`
	Data             [][]string `schema:"data"`
	ValueInputOption string     `schema:"valueInputOption"`
}

type GoogleSheetsAddDataToSpreadsheetExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[GoogleSheetsAddDataToSpreadsheetExecutorInput]
}

func NewGoogleSheetsAddDataToSpreadsheetExecutor(
	nodeID string,
	validator *jsonschema.Validator[GoogleSheetsAddDataToSpreadsheetExecutorInput],
) *GoogleSheetsAddDataToSpreadsheetExecutor {
	return &GoogleSheetsAddDataToSpreadsheetExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

type GoogleSheetsAddDataToSpreadsheetResponse struct {
	SpreadsheetID  string `json:"spreadsheetId"`
	UpdatedRange   string `json:"updatedRange"`
	UpdatedRows    int    `json:"updatedRows"`
	UpdatedColumns int    `json:"updatedColumns"`
	UpdatedCells   int    `json:"updatedCells"`
}

type GoogleSheetsAddDataToSpreadsheetErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func (e *GoogleSheetsAddDataToSpreadsheetExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			var successResponse GoogleSheetsAddDataToSpreadsheetResponse
			var errorResponse GoogleSheetsAddDataToSpreadsheetErrorResponse
			response, err := client.Post("/spreadsheets/"+input.SpreadsheetID+"/values/"+input.Range+":append?valueInputOption="+input.ValueInputOption).
				Body(map[string]any{
					"range":          input.Range,
					"majorDimension": "ROWS",
					"values":         input.Data,
				}).
				Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			results = append(results, map[string]any{
				"status":         true,
				"spreadsheetId":  successResponse.SpreadsheetID,
				"updatedRange":   successResponse.UpdatedRange,
				"updatedRows":    successResponse.UpdatedRows,
				"updatedColumns": successResponse.UpdatedColumns,
				"updatedCells":   successResponse.UpdatedCells,
			})
		}

		return results, nil
	}
}
