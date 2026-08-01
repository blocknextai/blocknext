package readspreadsheet

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/sheets/helpers"
)

type GoogleSheetsReadSpreadsheetExecutorInput struct {
	SpreadsheetID string `schema:"spreadsheetId"`
	Range         string `schema:"range"`
	FilterColumn  string `schema:"filterColumn"`
	FilterValue   string `schema:"filterValue"`
	Limit         int    `schema:"limit"`
	HeaderRow     int    `schema:"headerRow"`
}

type GoogleSheetsReadSpreadsheetExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[GoogleSheetsReadSpreadsheetExecutorInput]
}

func NewGoogleSheetsReadSpreadsheetExecutor(
	nodeID string,
	validator *jsonschema.Validator[GoogleSheetsReadSpreadsheetExecutorInput],
) *GoogleSheetsReadSpreadsheetExecutor {
	return &GoogleSheetsReadSpreadsheetExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

type GoogleSheetsReadSpreadsheetResponse struct {
	Range          string     `json:"range"`
	MajorDimension string     `json:"majorDimension"`
	Values         [][]string `json:"values"`
}

type GoogleSheetsReadSpreadsheetErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func (e *GoogleSheetsReadSpreadsheetExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			var response GoogleSheetsReadSpreadsheetResponse
			var errorResponse GoogleSheetsReadSpreadsheetErrorResponse

			res, err := client.Get("/spreadsheets/"+input.SpreadsheetID+"/values/"+input.Range).
				Do(&response, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !res.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			formattedData := helpers.FormatSpreadsheetData(response.Values, input.HeaderRow, input.FilterColumn, input.FilterValue, input.Limit)

			results = append(results, formattedData...)
		}

		return results, nil
	}
}
