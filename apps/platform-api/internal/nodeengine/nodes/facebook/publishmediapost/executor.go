package publishmediapost

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
	ErrEmptyMediaURL   = apperror.Internal("empty media url")
)

type FacebookPublishMediaPostExecutorInput struct {
	MediaURLs   []string `schema:"mediaUrls"`
	Message     string   `schema:"message"`
	AccountName string   `schema:"accountName"`
}

type FacebookPublishMediaPostExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[FacebookPublishMediaPostExecutorInput]
}

func NewFacebookPublishMediaPostExecutor(
	nodeID string,
	validator *jsonschema.Validator[FacebookPublishMediaPostExecutorInput],
) *FacebookPublishMediaPostExecutor {
	return &FacebookPublishMediaPostExecutor{
		ID:        nodeID,
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

type SuccessResponse struct {
	ID string `json:"id"`
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (e *FacebookPublishMediaPostExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			var postID string
			if len(input.MediaURLs) == 1 {
				mediaURL := strings.TrimSpace(input.MediaURLs[0])
				if mediaURL == "" {
					return nil, ErrEmptyMediaURL
				}

				mediaType := helpers.DetectMediaType(mediaURL)
				id, publishError := publishMediaPost(client, *account, mediaURL, input.Message, mediaType)
				if publishError != nil {
					return nil, publishError
				}
				postID = id
			} else {
				id, publishError := publishMultipleMediaPost(client, *account, input.MediaURLs, input.Message)
				if publishError != nil {
					return nil, publishError
				}
				postID = id
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

func publishMediaPost(
	client *httpclient.Client,
	account AccountData,
	mediaURL string,
	message string,
	mediaType string,
) (string, error) {
	var successResponse SuccessResponse
	var errorResponse ErrorResponse

	var endpoint string
	var bodyData map[string]any

	if mediaType == helpers.VideoType {
		endpoint = "/" + account.ID + "/videos"
		bodyData = map[string]any{
			"file_url":    mediaURL,
			"description": message,
		}
	} else {
		endpoint = "/" + account.ID + "/photos"
		bodyData = map[string]any{
			"url":     mediaURL,
			"caption": message,
		}
	}

	response, err := client.Post(endpoint).
		JSONContentType().
		QueryParam("access_token", account.AccessToken).
		Body(bodyData).
		Do(&successResponse, &errorResponse)

	if err != nil {
		return "", err
	}

	if !response.IsSuccess() {
		return "", apperror.Internal(errorResponse.Error.Message)
	}

	return successResponse.ID, nil
}

func publishMultipleMediaPost(
	client *httpclient.Client,
	account AccountData,
	mediaURLs []string,
	message string,
) (string, error) {
	var attachedMedia []map[string]string

	for _, mediaURL := range mediaURLs {
		mediaURL = strings.TrimSpace(mediaURL)
		if mediaURL == "" {
			continue
		}

		mediaType := helpers.DetectMediaType(mediaURL)

		mediaID, uploadErr := uploadMediaForAttachment(client, account, mediaURL, mediaType)
		if uploadErr != nil {
			return "", uploadErr
		}

		attachedMedia = append(attachedMedia, map[string]string{
			"media_fbid": mediaID,
		})
	}

	var successResponse SuccessResponse
	var errorResponse ErrorResponse

	response, err := client.Post("/"+account.ID+"/feed").
		JSONContentType().
		QueryParam("access_token", account.AccessToken).
		Body(map[string]any{
			"message":        message,
			"attached_media": attachedMedia,
		}).
		Do(&successResponse, &errorResponse)

	if err != nil {
		return "", err
	}

	if !response.IsSuccess() {
		return "", apperror.Internal(errorResponse.Error.Message)
	}

	return successResponse.ID, nil
}

func uploadMediaForAttachment(
	client *httpclient.Client,
	account AccountData,
	mediaURL string,
	mediaType string,
) (string, error) {
	var successResponse SuccessResponse
	var errorResponse ErrorResponse

	var endpoint string
	var bodyData map[string]any

	if mediaType == helpers.VideoType {
		endpoint = "/" + account.ID + "/videos"
		bodyData = map[string]any{
			"file_url":  mediaURL,
			"published": true,
		}
	} else {
		endpoint = "/" + account.ID + "/photos"
		bodyData = map[string]any{
			"url":       mediaURL,
			"published": true,
		}
	}

	response, err := client.Post(endpoint).
		JSONContentType().
		QueryParam("access_token", account.AccessToken).
		Body(bodyData).
		Do(&successResponse, &errorResponse)

	if err != nil {
		return "", err
	}

	if !response.IsSuccess() {
		return "", apperror.Internal(errorResponse.Error.Message)
	}

	return successResponse.ID, nil
}
