package sendmedia

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/whatsapp/helpers"
)

var (
	errMediaSendFailed = apperror.Internal("media send failed")
)

type WhatsAppSendMediaExecutorInput struct {
	PhoneNumber string   `schema:"phoneNumber"`
	MediaURLs   []string `schema:"mediaUrls"`
	Caption     string   `schema:"caption"`
}

type WhatsAppSendMediaExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[WhatsAppSendMediaExecutorInput]
}

func NewWhatsAppSendMediaExecutor(
	nodeID string,
	validator *jsonschema.Validator[WhatsAppSendMediaExecutorInput],
) *WhatsAppSendMediaExecutor {
	return &WhatsAppSendMediaExecutor{
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

func detectMediaType(url string) string {
	ext := strings.ToLower(filepath.Ext(url))

	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
		return "image"
	case ".mp4", ".3gp", ".mov":
		return "video"
	case ".mp3", ".aac", ".ogg", ".amr", ".m4a":
		return "audio"
	default:
		return "document"
	}
}

func (e *WhatsAppSendMediaExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			for _, mediaURL := range input.MediaURLs {
				messageID, err := e.sendSingleMediaByURL(client, phoneNumberID, input.PhoneNumber, mediaURL, input.Caption)
				if err != nil {
					results = append(results, map[string]any{
						"status": false,
						"error":  err.Error(),
					})
					continue
				}
				results = append(results, map[string]any{
					"status":    true,
					"messageId": messageID,
				})
			}
		}

		return results, nil
	}
}

func (e *WhatsAppSendMediaExecutor) sendSingleMediaByURL(client *httpclient.Client, phoneNumberID, to, mediaURL, caption string) (string, error) {
	mediaType := detectMediaType(mediaURL)

	mediaPayload := map[string]any{
		"link": mediaURL,
	}
	if caption != "" {
		mediaPayload["caption"] = caption
	}

	switch mediaType {
	case "document":
		mediaPayload["filename"] = filepath.Base(mediaURL)
	}

	var successResponse WhatsAppSuccessResponse
	var errorResponse WhatsAppErrorResponse

	response, err := client.Post("/"+phoneNumberID+"/messages").
		JSONContentType().
		Body(map[string]any{
			"messaging_product": "whatsapp",
			"recipient_type":    "individual",
			"to":                to,
			"type":              mediaType,
			mediaType:           mediaPayload,
		}).
		Do(&successResponse, &errorResponse)

	if err != nil {
		return "", err
	}

	if !response.IsSuccess() {
		return "", apperror.Internal(errorResponse.Error.Message)
	}

	if len(successResponse.Messages) == 0 {
		return "", errMediaSendFailed
	}

	return successResponse.Messages[0].ID, nil
}
