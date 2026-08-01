package searchemails

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/gmail/helpers"
)

type GmailSearchEmailsExecutorInput struct {
	Query string `schema:"query"`
}

type GmailSearchEmailsExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[GmailSearchEmailsExecutorInput]
}

type GmailSearchEmailsSuccessResponse struct {
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
	NextPageToken string `json:"nextPageToken"`
}

type GmailSearchEmailsErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

type GmailSearchEmailsResult struct {
	ID       string `json:"id"`
	ThreadID string `json:"threadId"`
	Snippet  string `json:"snippet"`
	Payload  struct {
		Headers []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"headers"`
		Body struct {
			Size int    `json:"size"`
			Data string `json:"data"`
		} `json:"body"`
	} `json:"payload"`
}

func NewGmailSearchEmailsExecutor(
	nodeID string,
	validator *jsonschema.Validator[GmailSearchEmailsExecutorInput],
) *GmailSearchEmailsExecutor {
	return &GmailSearchEmailsExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

func (e *GmailSearchEmailsExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "gmail_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")
		client := helpers.CreateClient(ctx, accessToken)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			emails, err := e.searchEmails(client, input.Query)
			if err != nil {
				return nil, err
			}

			results = append(results, map[string]any{
				"status": true,
				"result": map[string]any{
					"emails":      emails,
					"foundEmails": len(emails),
					"query":       input.Query,
				},
			})
		}

		return results, nil
	}
}

func (e *GmailSearchEmailsExecutor) searchEmails(client *httpclient.Client, query string) ([]map[string]any, error) {
	var searchResponse GmailSearchEmailsSuccessResponse
	var errorResponse GmailSearchEmailsErrorResponse

	response, err := client.Get("/users/me/messages").
		QueryParam("q", query).
		QueryParam("maxResults", "50").
		Do(&searchResponse, &errorResponse)

	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		return nil, apperror.Internal(errorResponse.Error.Message)
	}

	if len(searchResponse.Messages) == 0 {
		return []map[string]any{}, nil
	}

	emails := make([]map[string]any, 0, len(searchResponse.Messages))

	for _, message := range searchResponse.Messages {
		emailDetail, err := e.getEmailDetail(client, message.ID)
		if err != nil {
			continue
		}
		emails = append(emails, emailDetail)
	}

	return emails, nil
}

func (e *GmailSearchEmailsExecutor) getEmailDetail(client *httpclient.Client, emailID string) (map[string]any, error) {
	var emailResponse GmailSearchEmailsResult
	var errorResponse GmailSearchEmailsErrorResponse

	response, err := client.Get("/users/me/messages/"+emailID).
		QueryParam("format", "metadata").
		Do(&emailResponse, &errorResponse)

	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		return nil, apperror.Internal(errorResponse.Error.Message)
	}

	subject := ""
	sender := ""
	date := ""

	for _, header := range emailResponse.Payload.Headers {
		switch header.Name {
		case "Subject":
			subject = header.Value
		case "From":
			sender = header.Value
		case "Date":
			date = header.Value
		}
	}

	return map[string]any{
		"id":       emailResponse.ID,
		"threadId": emailResponse.ThreadID,
		"subject":  subject,
		"sender":   sender,
		"date":     date,
		"snippet":  emailResponse.Snippet,
	}, nil
}
