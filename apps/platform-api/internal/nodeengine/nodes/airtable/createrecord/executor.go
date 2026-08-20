package createrecord

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/airtable/helpers"
)

var (
	errExecutorRequestFailed = apperror.Internal("request failed")
	errExecutorInvalidFields = apperror.Internal("invalid fields")
)

type AirtableCreateRecordExecutorInput struct {
	BaseID  string `schema:"baseId"`
	TableID string `schema:"tableId"`
	Fields  string `schema:"fields"`
}

type AirtableCreateRecordExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[AirtableCreateRecordExecutorInput]
}

func NewAirtableCreateRecordExecutor(
	nodeID string,
	validator *jsonschema.Validator[AirtableCreateRecordExecutorInput],
) *AirtableCreateRecordExecutor {
	return &AirtableCreateRecordExecutor{
		ID:        nodeID,
		validator: validator,
	}
}

type SuccessResponse struct {
	ID          string         `json:"id"`
	CreatedTime string         `json:"createdTime"`
	Fields      map[string]any `json:"fields"`
}

type ErrorResponse struct {
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

func (e *AirtableCreateRecordExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "airtable_api")
		accessToken := credential.String("accessToken")
		client := helpers.CreateClient(ctx, accessToken)

		results := make([]map[string]any, 0, len(data))
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			var fields map[string]any
			if err := json.Unmarshal([]byte(input.Fields), &fields); err != nil {
				return nil, errExecutorInvalidFields
			}

			payload := map[string]any{
				"fields": fields,
			}

			var successResponse SuccessResponse
			var errorResponse ErrorResponse
			response, err := client.Post("/"+input.BaseID+"/"+input.TableID).
				JSONContentType().
				Body(payload).
				Do(&successResponse, &errorResponse)
			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				if errorResponse.Error.Message != "" {
					return nil, apperror.Internal(errorResponse.Error.Message)
				}
				return nil, errExecutorRequestFailed
			}

			results = append(results, map[string]any{
				"status": true,
				"id":     successResponse.ID,
			})
		}

		return results, nil
	}
}
