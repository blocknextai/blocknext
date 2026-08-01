package sendemail

import (
	"context"
	"mime"
	"net/mail"
	"strings"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/base64"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/gmail/helpers"
)

var (
	ErrExecutorEmptyResponse = apperror.Internal("empty response")
	ErrInvalidRecipient      = apperror.Internal("invalid recipient")
	ErrFailedToResolveSender = apperror.Internal("failed to resolve gmail sender")
)

type GmailSendEmailExecutorInput struct {
	To      string `schema:"to"`
	Subject string `schema:"subject"`
	Body    string `schema:"body"`
}

type GmailSendEmailExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[GmailSendEmailExecutorInput]
}

func NewGmailSendEmailExecutor(
	nodeID string,
	validator *jsonschema.Validator[GmailSendEmailExecutorInput],
) *GmailSendEmailExecutor {
	return &GmailSendEmailExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

type GmailProfileResponse struct {
	EmailAddress string `json:"emailAddress"`
}

type GmailSendEmailSuccessResponse struct {
	ID       string   `json:"id"`
	ThreadID string   `json:"threadId"`
	LabelIDs []string `json:"labelIds"`
}

type GmailSendEmailErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func (e *GmailSendEmailExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "gmail_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")
		client := helpers.CreateClient(ctx, accessToken)

		var profileResponse GmailProfileResponse
		var profileErrorResponse GmailSendEmailErrorResponse
		profileRes, err := client.Get("/users/me/profile").
			Do(&profileResponse, &profileErrorResponse)
		if err != nil {
			return nil, err
		}
		if !profileRes.IsSuccess() || strings.TrimSpace(profileResponse.EmailAddress) == "" {
			return nil, ErrFailedToResolveSender
		}
		fromAddr := (new(mail.Address{Address: profileResponse.EmailAddress})).String()

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			toAddr, err := mail.ParseAddress(input.To)
			if err != nil {
				return nil, ErrInvalidRecipient
			}

			raw := []byte(
				"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
					"From: " + fromAddr + "\r\n" +
					"To: " + toAddr.String() + "\r\n" +
					"Subject: " + mime.QEncoding.Encode("UTF-8", input.Subject) + "\r\n\r\n" +
					input.Body,
			)

			encoded := base64.URLEncode(raw)

			var successResponse GmailSendEmailSuccessResponse
			var errorResponse GmailSendEmailErrorResponse
			response, err := client.Post("/users/me/messages/send").
				JSONContentType().
				Body(map[string]any{
					"raw": encoded,
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
				"status": true,
			})
		}

		return results, nil
	}
}
