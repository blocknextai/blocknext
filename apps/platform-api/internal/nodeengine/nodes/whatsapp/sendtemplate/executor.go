package sendtemplate

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/whatsapp/helpers"
)

var (
	errTemplateSendFailed = apperror.Internal("template send failed")
)

type WhatsAppSendTemplateExecutorInput struct {
	PhoneNumber  string `schema:"phoneNumber"`
	TemplateName string `schema:"templateName"`
	LanguageCode string `schema:"languageCode"`
}

type WhatsAppSendTemplateExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[WhatsAppSendTemplateExecutorInput]
}

func NewWhatsAppSendTemplateExecutor(
	nodeID string,
	validator *jsonschema.Validator[WhatsAppSendTemplateExecutorInput],
) *WhatsAppSendTemplateExecutor {
	return &WhatsAppSendTemplateExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

type WhatsAppSuccessResponse struct {
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
}

type WhatsAppErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (e *WhatsAppSendTemplateExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "whatsapp_api")
		accessToken := credential.String("accessToken")
		phoneNumberID := credential.String("phoneNumberId")
		client := helpers.CreateClient(ctx, accessToken)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			var successResponse WhatsAppSuccessResponse
			var errorResponse WhatsAppErrorResponse

			response, err := client.Post("/"+phoneNumberID+"/messages").
				JSONContentType().
				Body(map[string]any{
					"messaging_product": "whatsapp",
					"to":                input.PhoneNumber,
					"type":              "template",
					"template": map[string]any{
						"name": input.TemplateName,
						"language": map[string]any{
							"code": input.LanguageCode,
						},
					},
				}).
				Do(&successResponse, &errorResponse)

			if err != nil {
				return nil, err
			}

			if !response.IsSuccess() {
				return nil, apperror.Internal(errorResponse.Error.Message)
			}

			if len(successResponse.Messages) == 0 {
				return nil, errTemplateSendFailed
			}

			messageID := ""
			if len(successResponse.Messages) > 0 {
				messageID = successResponse.Messages[0].ID
			}

			results = append(results, map[string]any{
				"status":    true,
				"messageId": messageID,
			})
		}

		return results, nil
	}
}
