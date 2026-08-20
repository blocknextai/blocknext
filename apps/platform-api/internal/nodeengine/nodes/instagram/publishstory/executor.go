package publishstory

import (
	"context"
	"strings"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/instagram/helpers"
)

var (
	ErrEmptyPublishID = apperror.Internal("instagram returned empty story id")
)

type InstagramPublishStoryExecutorInput struct {
	MediaURLs []string `schema:"mediaUrls"`
	StoryLink string   `schema:"storyLink"`
}

type InstagramPublishStoryExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[InstagramPublishStoryExecutorInput]
}

func NewInstagramPublishStoryExecutor(
	nodeID string,
	validator *jsonschema.Validator[InstagramPublishStoryExecutorInput],
) *InstagramPublishStoryExecutor {
	return &InstagramPublishStoryExecutor{
		ID:        nodeID,
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

func (e *InstagramPublishStoryExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			for _, mediaURL := range input.MediaURLs {
				mediaURL = strings.TrimSpace(mediaURL)
				if mediaURL == "" {
					continue
				}

				detectedType := helpers.DetectMediaType(mediaURL)

				mediaID, mediaError := createStoryMedia(ctx, accessToken, mediaURL, detectedType, input.StoryLink)
				if mediaError != nil {
					results = append(results, map[string]any{
						"status": false,
						"error":  mediaError.Error(),
					})
					continue
				}

				if detectedType == helpers.VideoType {
					waitError := helpers.WaitForMediaReady(ctx, accessToken, mediaID)
					if waitError != nil {
						results = append(results, map[string]any{
							"status": false,
							"error":  waitError.Error(),
						})
						continue
					}
				}

				publishedID, publishError := publishStory(ctx, accessToken, mediaID)
				if publishError != nil {
					results = append(results, map[string]any{
						"status": false,
						"error":  publishError.Error(),
					})
					continue
				}

				if publishedID == "" {
					results = append(results, map[string]any{
						"status": false,
						"error":  ErrEmptyPublishID.Error(),
					})
					continue
				}

				results = append(results, map[string]any{
					"status":  true,
					"storyId": publishedID,
				})
			}
		}

		return results, nil
	}
}

func createStoryMedia(
	ctx context.Context,
	accessToken string,
	mediaURL string,
	mediaType string,
	storyLink string,
) (string, error) {
	client := helpers.GetInstagramClient(ctx, accessToken)

	mediaRequestBuilder := client.Post("/me/media").
		JSONContentType()

	if mediaType == helpers.VideoType {
		mediaRequestBuilder.QueryParam("video_url", mediaURL)
	} else {
		mediaRequestBuilder.QueryParam("image_url", mediaURL)
	}

	mediaRequestBuilder.QueryParam("media_type", helpers.StoriesType)

	if storyLink != "" {
		mediaRequestBuilder.QueryParam("story_link", storyLink)
	}

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

func publishStory(
	ctx context.Context,
	accessToken string,
	mediaID string,
) (string, error) {
	client := helpers.GetInstagramClient(ctx, accessToken)

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
