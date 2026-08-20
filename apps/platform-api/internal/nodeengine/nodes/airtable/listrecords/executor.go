package listrecords

import (
	"context"
	"strconv"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/airtable/helpers"
)

var (
	errExecutorRequestFailed = apperror.Internal("request failed")
)

type AirtableListRecordsExecutorInput struct {
	BaseID     string `schema:"baseId"`
	TableID    string `schema:"tableId"`
	MaxRecords int    `schema:"maxRecords"`
}

type AirtableListRecordsExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[AirtableListRecordsExecutorInput]
}

func NewAirtableListRecordsExecutor(
	nodeID string,
	validator *jsonschema.Validator[AirtableListRecordsExecutorInput],
) *AirtableListRecordsExecutor {
	return &AirtableListRecordsExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

type Record struct {
	ID          string         `json:"id"`
	CreatedTime string         `json:"createdTime"`
	Fields      map[string]any `json:"fields"`
}

type SuccessResponse struct {
	Records []Record `json:"records"`
	Offset  string   `json:"offset"`
}

type ErrorResponse struct {
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

func (e *AirtableListRecordsExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "airtable_api")
		accessToken := credential.String("accessToken")
		client := helpers.CreateClient(ctx, accessToken)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			request := client.Get("/" + input.BaseID + "/" + input.TableID).
				JSONContentType()
			if input.MaxRecords > 0 {
				request = request.QueryParam("maxRecords", strconv.Itoa(input.MaxRecords))
			}

			var successResponse SuccessResponse
			var errorResponse ErrorResponse
			response, err := request.Do(&successResponse, &errorResponse)
			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				if errorResponse.Error.Message != "" {
					return nil, apperror.Internal(errorResponse.Error.Message)
				}
				return nil, errExecutorRequestFailed
			}

			for _, record := range successResponse.Records {
				results = append(results, map[string]any{
					"id":          record.ID,
					"createdTime": record.CreatedTime,
					"fields":      record.Fields,
				})
			}
		}

		return results, nil
	}
}
