package publishstory

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

const (
	VideoStatusMaxRetry = 30
)

var (
	ErrAccountNotFound = apperror.Internal("account not found")
	ErrEmptyPublishID  = apperror.Internal("facebook returned empty story id")
)

type FacebookPublishStoryExecutorInput struct {
	MediaURLs   []string `schema:"mediaUrls"`
	AccountName string   `schema:"accountName"`
}

type FacebookPublishStoryExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[FacebookPublishStoryExecutorInput]
}

func NewFacebookPublishStoryExecutor(
	nodeID string,
	validator *jsonschema.Validator[FacebookPublishStoryExecutorInput],
) *FacebookPublishStoryExecutor {
	return &FacebookPublishStoryExecutor{
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

type SuccessResponse struct {
	ID string `json:"id"`
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (e *FacebookPublishStoryExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			for _, mediaURL := range input.MediaURLs {
				mediaURL = strings.TrimSpace(mediaURL)
				if mediaURL == "" {
					continue
				}

				mediaType := helpers.DetectMediaType(mediaURL)
				if mediaType == helpers.VideoType {
					continue
				}

				id, publishError := publishStory(ctx, client, *account, mediaURL, mediaType)
				if publishError != nil {
					results = append(results, map[string]any{
						"status": false,
						"error":  publishError.Error(),
					})
					continue
				}

				if id == "" {
					results = append(results, map[string]any{
						"status": false,
						"error":  ErrEmptyPublishID.Error(),
					})
					continue
				}

				results = append(results, map[string]any{
					"status":   true,
					"storyId":  id,
					"storyUrl": "https://www.facebook.com/" + id,
				})
			}
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

func publishStory(
	ctx context.Context,
	client *httpclient.Client,
	account AccountData,
	mediaURL string,
	mediaType string,
) (string, error) {
	/*if mediaType == helpers.VideoType {
		return publishVideoStory(ctx, client, account, mediaURL)
	}*/

	if mediaType == helpers.PhotoType {
		return publishPhotoStory(ctx, client, account, mediaURL)
	}

	return "", nil
}

type StoryResponse struct {
	Success bool   `json:"success"`
	PostID  string `json:"post_id"`
}

func publishPhotoStory(
	ctx context.Context,
	client *httpclient.Client,
	account AccountData,
	mediaURL string,
) (string, error) {
	var uploadResponse SuccessResponse
	var errorResponse ErrorResponse

	uploadResp, uploadErr := client.Post("/"+account.ID+"/photos").
		JSONContentType().
		QueryParam("access_token", account.AccessToken).
		Body(map[string]any{
			"url":       mediaURL,
			"published": false,
		}).
		Do(&uploadResponse, &errorResponse)

	if uploadErr != nil {
		return "", uploadErr
	}

	if !uploadResp.IsSuccess() {
		return "", apperror.Internal(errorResponse.Error.Message)
	}

	photoID := uploadResponse.ID

	var storyResponse StoryResponse
	storyResp, storyErr := client.Post("/"+account.ID+"/photo_stories").
		JSONContentType().
		QueryParam("access_token", account.AccessToken).
		Body(map[string]any{
			"photo_id": photoID,
		}).
		Do(&storyResponse, &errorResponse)

	if storyErr != nil {
		return "", storyErr
	}

	if !storyResp.IsSuccess() {
		return "", apperror.Internal(errorResponse.Error.Message)
	}

	return storyResponse.PostID, nil
}

type VideoInitResponse struct {
	VideoID   string `json:"video_id"`
	UploadURL string `json:"upload_url"`
}

type UploadResponse struct {
	Success bool `json:"success"`
}

type VideoStatusResponse struct {
	Status struct {
		VideoStatus    string `json:"video_status"`
		UploadingPhase struct {
			Status          string `json:"status"`
			BytesTransfered int64  `json:"bytes_transfered"`
		} `json:"uploading_phase"`
		ProcessingPhase struct {
			Status string `json:"status"`
			Error  *struct {
				Message string `json:"message"`
			} `json:"error,omitempty"`
		} `json:"processing_phase"`
		PublishingPhase struct {
			Status        string `json:"status"`
			PublishStatus string `json:"publish_status"`
			PublishTime   int64  `json:"publish_time"`
		} `json:"publishing_phase"`
	} `json:"status"`
}
