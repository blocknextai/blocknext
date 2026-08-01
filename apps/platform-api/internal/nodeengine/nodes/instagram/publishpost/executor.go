package publishpost

import (
	"context"
	"strings"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/instagram/helpers"
)

var (
	errNoMediaURLsProvided = apperror.Internal("no media urls provided")
)

type InstagramPublishPostExecutorInput struct {
	MediaURLs []string `schema:"mediaUrls"`
	Caption   string   `schema:"caption"`
	Location  string   `schema:"location"`
	AltText   string   `schema:"altText"`
}

type InstagramPublishPostExecutor struct {
	executors.Executor
	validator *jsonschema.Validator[InstagramPublishPostExecutorInput]
}

func NewInstagramPublishPostExecutor(
	nodeID string,
	validator *jsonschema.Validator[InstagramPublishPostExecutorInput],
) *InstagramPublishPostExecutor {
	return &InstagramPublishPostExecutor{
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

func (e *InstagramPublishPostExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
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

			var publishedMediaID string
			if len(input.MediaURLs) > 1 {
				var mediaIDs []string
				for _, mediaURL := range input.MediaURLs {
					detectedType := helpers.DetectMediaType(mediaURL)
					mediaID, mediaError := createCarouselMediaWithType(client, mediaURL, detectedType, input.AltText)
					if mediaError != nil {
						return nil, mediaError
					}

					if detectedType == helpers.VideoType {
						waitError := helpers.WaitForMediaReady(ctx, accessToken, mediaID)
						if waitError != nil {
							return nil, waitError
						}
					}

					mediaIDs = append(mediaIDs, mediaID)
				}
				containerID, containerError := createCarouselContainer(client, mediaIDs, input.Caption, input.Location)
				if containerError != nil {
					return nil, containerError
				}
				postID, publishError := publishPost(client, containerID)
				if publishError != nil {
					return nil, publishError
				}
				publishedMediaID = postID
			} else if len(input.MediaURLs) == 1 {
				mediaURL := input.MediaURLs[0]
				detectedType := helpers.DetectMediaType(mediaURL)
				mediaID, mediaError := createMedia(client, mediaURL, input.Caption, detectedType, input.Location, input.AltText)
				if mediaError != nil {
					return nil, mediaError
				}

				if detectedType == helpers.VideoType {
					waitError := helpers.WaitForMediaReady(ctx, accessToken, mediaID)
					if waitError != nil {
						return nil, waitError
					}
				}

				postID, publishError := publishPost(client, mediaID)
				if publishError != nil {
					return nil, publishError
				}
				publishedMediaID = postID
			} else {
				return nil, errNoMediaURLsProvided
			}

			results = append(results, map[string]any{
				"status": true,
				"postId": publishedMediaID,
			})
		}

		return results, nil
	}
}

func createMedia(
	client *httpclient.Client,
	mediaURL string,
	caption string,
	mediaType string,
	location string,
	altText string,
) (string, error) {
	mediaRequestBuilder := client.Post("/me/media").
		JSONContentType()

	if mediaType == helpers.VideoType {
		mediaRequestBuilder.QueryParam("video_url", mediaURL)
		mediaRequestBuilder.QueryParam("media_type", helpers.VideoType)
	} else {
		mediaRequestBuilder.QueryParam("image_url", mediaURL)
	}

	if caption != "" {
		mediaRequestBuilder.QueryParam("caption", caption)
	}

	if location != "" {
		mediaRequestBuilder.QueryParam("location_id", location)
	}

	if altText != "" {
		mediaRequestBuilder.QueryParam("alt_text", altText)
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

func createCarouselMediaWithType(
	client *httpclient.Client,
	mediaURL string,
	mediaType string,
	altText string,
) (string, error) {
	mediaRequestBuilder := client.Post("/me/media").
		JSONContentType().
		QueryParam("is_carousel_item", "true")

	if mediaType == helpers.VideoType {
		mediaRequestBuilder.QueryParam("video_url", mediaURL)
		mediaRequestBuilder.QueryParam("media_type", helpers.VideoType)
	} else {
		mediaRequestBuilder.QueryParam("image_url", mediaURL)
	}

	if altText != "" {
		mediaRequestBuilder.QueryParam("alt_text", altText)
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

func createCarouselContainer(
	client *httpclient.Client,
	mediaIDs []string,
	caption string,
	location string,
) (string, error) {
	containerRequestBuilder := client.Post("/me/media").
		JSONContentType().
		QueryParam("media_type", helpers.CarouselType)

	childrenParam := strings.Join(mediaIDs, ",")
	containerRequestBuilder.QueryParam("children", childrenParam)

	if caption != "" {
		containerRequestBuilder.QueryParam("caption", caption)
	}

	if location != "" {
		containerRequestBuilder.QueryParam("location_id", location)
	}

	var containerSuccessResponse SuccessResponse
	var containerErrorResponse ErrorResponse
	containerResponse, err := containerRequestBuilder.Do(&containerSuccessResponse, &containerErrorResponse)

	if err != nil {
		return "", err
	}

	if !containerResponse.IsSuccess() {
		return "", apperror.Internal(containerErrorResponse.Error.Message)
	}

	return containerSuccessResponse.ID, nil
}

func publishPost(
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
