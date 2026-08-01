package publishreels

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/instagram/helpers"
)

type InstagramPublishReelsExecutorInput struct {
	VideoURL      string `schema:"videoUrl"`
	Caption       string `schema:"caption"`
	CoverImageURL string `schema:"coverImageUrl"`
	ShareToFeed   bool   `schema:"shareToFeed"`
}

type InstagramPublishReelsExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[InstagramPublishReelsExecutorInput]
}

func NewInstagramPublishReelsExecutor(
	nodeID string,
	validator *jsonschema.Validator[InstagramPublishReelsExecutorInput],
) *InstagramPublishReelsExecutor {
	return &InstagramPublishReelsExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator: validator,
	}
}

type SuccessResponse struct {
	ID string `json:"id"`
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (e *InstagramPublishReelsExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "instagram_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			client := helpers.GetInstagramClient(ctx, accessToken)

			mediaID, mediaError := createReelsMedia(client, input.VideoURL, input.Caption, input.CoverImageURL, input.ShareToFeed)
			if mediaError != nil {
				return nil, mediaError
			}

			waitError := helpers.WaitForMediaReady(ctx, accessToken, mediaID)
			if waitError != nil {
				return nil, waitError
			}

			reelsID, publishError := publishReels(client, mediaID)
			if publishError != nil {
				return nil, publishError
			}

			results = append(results, map[string]any{
				"status":  true,
				"reelsId": reelsID,
			})
		}

		return results, nil
	}
}

func createReelsMedia(
	client *httpclient.Client,
	videoURL string,
	caption string,
	coverImageURL string,
	shareToFeed bool,
) (string, error) {
	mediaRequestBuilder := client.Post("/me/media").
		JSONContentType().
		QueryParam("video_url", videoURL).
		QueryParam("media_type", helpers.ReelsType)

	if caption != "" {
		mediaRequestBuilder.QueryParam("caption", caption)
	}

	if coverImageURL != "" {
		mediaRequestBuilder.QueryParam("thumb_offset", "0").
			QueryParam("cover_url", coverImageURL)
	}

	mediaRequestBuilder.QueryParam("share_to_feed", boolToString(shareToFeed))

	var mediaSuccessResponse SuccessResponse
	var mediaErrorResponse ErrorResponse
	mediaResponse, err := mediaRequestBuilder.Do(&mediaSuccessResponse, &mediaErrorResponse)

	if err != nil {
		return "", err
	}

	if !mediaResponse.IsSuccess() {
		return "", apperror.Internal(mediaErrorResponse.Error.Message)
	}

	return mediaSuccessResponse.ID, nil
}

func publishReels(
	client *httpclient.Client,
	mediaID string,
) (string, error) {
	var mediaPublishSuccessResponse SuccessResponse
	var mediaPublishErrorResponse ErrorResponse
	mediaPublishResponse, err := client.Post("/me/media_publish").
		JSONContentType().
		QueryParam("creation_id", mediaID).
		Do(&mediaPublishSuccessResponse, &mediaPublishErrorResponse)

	if err != nil {
		return "", err
	}

	if !mediaPublishResponse.IsSuccess() {
		return "", apperror.Internal(mediaPublishErrorResponse.Error.Message)
	}

	return mediaPublishSuccessResponse.ID, nil
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
