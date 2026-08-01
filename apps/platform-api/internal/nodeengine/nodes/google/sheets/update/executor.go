package update

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/cast"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/sheets/helpers"
)

var (
	errColumnNotFound = apperror.Internal("column not found")
	errInvalidColumn  = apperror.Internal("invalid column")
)

type GoogleSheetsUpdateExecutorInput struct {
	SheetID   string `schema:"sheetId"`
	RowNumber int    `schema:"rowNumber"`
	Column    string `schema:"column"`
	NewValue  string `schema:"newValue"`
}

type GoogleSheetsUpdateExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[GoogleSheetsUpdateExecutorInput]
}

func NewGoogleSheetsUpdateExecutor(
	nodeID string,
	validator *jsonschema.Validator[GoogleSheetsUpdateExecutorInput],
) *GoogleSheetsUpdateExecutor {
	return &GoogleSheetsUpdateExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

type GoogleSheetsResponse struct {
	SpreadsheetID  string `json:"spreadsheetId"`
	UpdatedRange   string `json:"updatedRange"`
	UpdatedRows    int    `json:"updatedRows"`
	UpdatedColumns int    `json:"updatedColumns"`
	UpdatedCells   int    `json:"updatedCells"`
}

type GoogleSheetsReadResponse struct {
	Range          string     `json:"range"`
	MajorDimension string     `json:"majorDimension"`
	Values         [][]string `json:"values"`
}

type GoogleSheetsErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func (e *GoogleSheetsUpdateExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			var readResponse GoogleSheetsReadResponse
			var errorResponse GoogleSheetsErrorResponse

			headerRange := "1:1"
			response, err := client.Get("/spreadsheets/"+input.SheetID+"/values/"+headerRange).
				Do(&readResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			var columnLetter string
			if len(readResponse.Values) > 0 && len(readResponse.Values[0]) > 0 {
				headerRow := readResponse.Values[0]
				columnIndex, err := helpers.FindColumnIndex(input.Column, headerRow)
				if err != nil {
					return nil, errColumnNotFound
				}
				columnLetter = helpers.IndexToColumnLetter(columnIndex)
			} else {
				if _, err := helpers.ColumnLetterToIndex(input.Column); err == nil {
					columnLetter = input.Column
				} else {
					return nil, errInvalidColumn
				}
			}

			cellAddress := columnLetter + cast.ToString(input.RowNumber)

			var successResponse GoogleSheetsResponse
			response, err = client.Put("/spreadsheets/"+input.SheetID+"/values/"+cellAddress+"?valueInputOption=USER_ENTERED").
				Body(map[string]any{
					"range":          cellAddress,
					"majorDimension": "ROWS",
					"values":         [][]string{{input.NewValue}},
				}).
				Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			results = append(results, map[string]any{
				"status":       true,
				"sheetId":      successResponse.SpreadsheetID,
				"updatedRange": successResponse.UpdatedRange,
				"rowNumber":    input.RowNumber,
				"column":       input.Column,
				"newValue":     input.NewValue,
				"updatedCells": successResponse.UpdatedCells,
			})
		}

		return results, nil
	}
}
