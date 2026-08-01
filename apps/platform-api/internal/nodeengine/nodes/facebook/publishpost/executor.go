package publishpost

import (
	"context"
	"strings"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/facebook/helpers"
)

var (
	ErrAccountNotFound = apperror.Internal("account not found")
)

type FacebookPublishPostExecutorInput struct {
	AccountName string `schema:"accountName"`
	Message     string `schema:"message"`
}

type FacebookPublishPostExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[FacebookPublishPostExecutorInput]
}

func NewFacebookPublishPostExecutor(
	nodeID string,
	validator *jsonschema.Validator[FacebookPublishPostExecutorInput],
) *FacebookPublishPostExecutor {
	return &FacebookPublishPostExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

type AccountData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AccessToken string `json:"access_token"`
}

type SuccessAccountsResponse struct {
	Data []AccountData `json:"data"`
}

type SuccessFeedResponse struct {
	ID string `json:"id"`
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (e *FacebookPublishPostExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "facebook_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			client := helpers.GetFacebookClient(ctx, accessToken)

			account, accountError := getAccount(client, input.AccountName)
			if accountError != nil {
				return nil, accountError
			}

			postID, feedError := publishFeed(client, *account, input.Message)
			if feedError != nil {
				return nil, feedError
			}

			postURL := ""
			if postID != "" {
				postURL = "https://www.facebook.com/" + postID
			}

			results = append(results, map[string]any{
				"status":  true,
				"postId":  postID,
				"postUrl": postURL,
			})
		}

		return results, nil
	}
}

func getAccount(client *httpclient.Client, accountName string) (*AccountData, error) {
	var successAccountsResponse SuccessAccountsResponse
	var errorAccountsResponse ErrorResponse
	accountsResponse, err := client.Get("/me/accounts").
		JSONContentType().
		Do(&successAccountsResponse, &errorAccountsResponse)

	if err != nil {
		return nil, err
	}

	if !accountsResponse.IsSuccess() {
		return nil, apperror.Internal(errorAccountsResponse.Error.Message)
	}

	if len(successAccountsResponse.Data) == 0 {
		return nil, ErrAccountNotFound
	}

	accountData := successAccountsResponse.Data[0]
	if accountName != "" {
		for _, account := range successAccountsResponse.Data {
			if strings.EqualFold(account.Name, accountName) {
				accountData = account
				break
			}
		}
	}

	if accountData.ID == "" {
		return nil, ErrAccountNotFound
	}

	return &accountData, nil
}

func publishFeed(
	client *httpclient.Client,
	account AccountData,
	message string,
) (string, error) {
	var successFeedResponse SuccessFeedResponse
	var errorFeedResponse ErrorResponse
	feedResponse, err := client.Post("/"+account.ID+"/feed").
		JSONContentType().
		QueryParam("access_token", account.AccessToken).
		Body(map[string]any{
			"message": message,
		}).
		Do(&successFeedResponse, &errorFeedResponse)

	if err != nil {
		return "", err
	}

	if !feedResponse.IsSuccess() {
		return "", apperror.Internal(errorFeedResponse.Error.Message)
	}

	return successFeedResponse.ID, nil
}
