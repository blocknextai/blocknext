package publishmediapost

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/x/helpers"
)

const (
	MediaUploadMaxRetry = 30
)

var (
	ErrFailedToInitializeUpload   = apperror.Internal("failed to initialize upload")
	ErrFailedToAppendMedia        = apperror.Internal("failed to append media")
	ErrFailedToFinalizeUpload     = apperror.Internal("failed to finalize upload")
	ErrMediaUploadMaxRetryReached = apperror.Internal("media upload max retry reached")
)

type XPublishMediaPostExecutorInput struct {
	MediaURLs      []string `schema:"mediaUrls"`
	Text           string   `schema:"text"`
	ReplyToTweetID string   `schema:"replyToTweetId"`
}

type XPublishMediaPostExecutor struct {
	executors.Executor
	validator   *jsonschema.Validator[XPublishMediaPostExecutorInput]
	fileGateway filegateway.FileGateway
}

func NewXPublishMediaPostExecutor(
	nodeID string,
	validator *jsonschema.Validator[XPublishMediaPostExecutorInput],
	fileGateway filegateway.FileGateway,
) *XPublishMediaPostExecutor {
	return &XPublishMediaPostExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator:   validator,
		fileGateway: fileGateway,
	}
}

func (e *XPublishMediaPostExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "x_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			var mediaIDs []string
			var rateLimitInfo helpers.RateLimitInfo
			for _, mediaURL := range input.MediaURLs {
				mediaURL = strings.TrimSpace(mediaURL)
				if mediaURL == "" {
					continue
				}

				mediaID, rateLimit, uploadErr := e.uploadMedia(ctx, accessToken, mediaURL)
				rateLimitInfo = rateLimit
				if uploadErr != nil {
					result := map[string]any{
						"status":    false,
						"rateLimit": helpers.RateLimitInfoToMap(rateLimitInfo),
					}
					results = append(results, result)
					return results, uploadErr
				}

				mediaIDs = append(mediaIDs, mediaID)
			}

			rateLimitInfo, err = publishTweetWithMedia(ctx, accessToken, input.Text, mediaIDs, input.ReplyToTweetID)

			result := map[string]any{
				"rateLimit": helpers.RateLimitInfoToMap(rateLimitInfo),
			}

			if err != nil {
				result["status"] = false
				results = append(results, result)
				return results, err
			}

			result["status"] = true
			results = append(results, result)
		}

		return results, nil
	}
}

func (e *XPublishMediaPostExecutor) uploadMedia(ctx context.Context, accessToken string, mediaURL string) (string, helpers.RateLimitInfo, error) {
	var rateLimitInfo helpers.RateLimitInfo

	downloadResult, err := e.fileGateway.DownloadFile(ctx, mediaURL)
	if err != nil {
		return "", rateLimitInfo, err
	}
	defer func() {
		_ = downloadResult.DataReader.Close()
	}()

	mediaData, err := io.ReadAll(downloadResult.DataReader)
	if err != nil {
		return "", rateLimitInfo, err
	}

	mediaCategory := helpers.DetectMediaCategoryFromContentType(downloadResult.ContentType, mediaURL)

	mediaID, rateLimit, err := initializeMediaUpload(ctx, accessToken, len(mediaData), mediaCategory)
	rateLimitInfo = rateLimit
	if err != nil {
		return "", rateLimitInfo, err
	}

	rateLimit, err = appendMediaData(ctx, accessToken, mediaID, mediaData)
	rateLimitInfo = rateLimit
	if err != nil {
		return "", rateLimitInfo, err
	}

	mediaKey, rateLimit, err := finalizeMediaUpload(ctx, accessToken, mediaID)
	rateLimitInfo = rateLimit
	if err != nil {
		return "", rateLimitInfo, err
	}

	return mediaKey, rateLimitInfo, nil
}

type MediaUploadInitResponse struct {
	Data struct {
		ID       string `json:"id"`
		MediaKey string `json:"media_key"`
	} `json:"data"`
}

func initializeMediaUpload(ctx context.Context, accessToken string, totalBytes int, mediaCategory string) (string, helpers.RateLimitInfo, error) {
	client := helpers.GetXClient(ctx, accessToken)

	reqBody := map[string]any{
		"total_bytes":    totalBytes,
		"media_category": mediaCategory,
	}

	var successResponse MediaUploadInitResponse
	var errorResponse helpers.ErrorResponse
	response, err := client.Post("/2/media/upload/initialize").
		JSONContentType().
		Body(reqBody).
		Do(&successResponse, &errorResponse)

	rateLimitInfo := helpers.ExtractRateLimitInfo(response.Headers)

	if err != nil {
		return "", rateLimitInfo, err
	}

	if !response.IsSuccess() {
		return "", rateLimitInfo, apperror.Internal(errorResponse.Detail)
	}

	if successResponse.Data.MediaKey == "None" {
		return "", rateLimitInfo, ErrFailedToInitializeUpload
	}

	return successResponse.Data.ID, rateLimitInfo, nil
}

func appendMediaData(ctx context.Context, accessToken string, mediaID string, mediaData []byte) (helpers.RateLimitInfo, error) {
	client := helpers.GetXClient(ctx, accessToken)

	var errorResponse helpers.ErrorResponse
	response, err := client.Post("/2/media/upload/"+mediaID+"/append").
		MultipartFormContentType().
		MultipartFormBody().
		AddFile("media", "media", mediaData).
		AddField("segment_index", "0").
		Do(nil, &errorResponse)

	rateLimitInfo := helpers.ExtractRateLimitInfo(response.Headers)

	if err != nil {
		return rateLimitInfo, ErrFailedToAppendMedia
	}

	return rateLimitInfo, nil
}

type MediaUploadFinalizeResponse struct {
	Data struct {
		ID             string `json:"id"`
		MediaKey       string `json:"media_key"`
		ProcessingInfo *struct {
			State          string `json:"state"`
			CheckAfterSecs int    `json:"check_after_secs"`
		} `json:"processing_info"`
	} `json:"data"`
}

func finalizeMediaUpload(ctx context.Context, accessToken string, mediaID string) (string, helpers.RateLimitInfo, error) {
	client := helpers.GetXClient(ctx, accessToken)

	var successResponse MediaUploadFinalizeResponse
	var errorResponse helpers.ErrorResponse
	response, err := client.Post("/2/media/upload/"+mediaID+"/finalize").
		JSONContentType().
		Do(&successResponse, &errorResponse)

	rateLimitInfo := helpers.ExtractRateLimitInfo(response.Headers)

	if err != nil {
		return "", rateLimitInfo, err
	}

	if !response.IsSuccess() {
		return "", rateLimitInfo, apperror.Internal(errorResponse.Detail)
	}

	info := successResponse.Data.ProcessingInfo
	if info != nil {
		if info.State == helpers.ProcessingInfoStateFailed {
			return "", rateLimitInfo, ErrFailedToFinalizeUpload
		}

		if info.State == helpers.ProcessingInfoStatePending || info.State == helpers.ProcessingInfoStateInProgress {
			rateLimit, err := waitForMediaProcessing(ctx, accessToken, mediaID)
			rateLimitInfo = rateLimit
			if err != nil {
				return "", rateLimitInfo, err
			}
		}
	}

	return successResponse.Data.ID, rateLimitInfo, nil
}

func waitForMediaProcessing(ctx context.Context, accessToken string, mediaID string) (helpers.RateLimitInfo, error) {
	client := helpers.GetXClient(ctx, accessToken)
	retryDelayDuration := time.Duration(20_000) * time.Millisecond
	retryCount := 0
	var rateLimitInfo helpers.RateLimitInfo

	for {
		select {
		case <-ctx.Done():
			return rateLimitInfo, ctx.Err()
		case <-time.After(retryDelayDuration):
		}

		if retryCount++; retryCount > MediaUploadMaxRetry {
			return rateLimitInfo, ErrMediaUploadMaxRetryReached
		}

		var successResponse MediaUploadFinalizeResponse
		var errorResponse helpers.ErrorResponse
		response, err := client.Get("/2/media/upload/").
			JSONContentType().
			QueryParam("media_id", mediaID).
			Do(&successResponse, &errorResponse)

		rateLimitInfo = helpers.ExtractRateLimitInfo(response.Headers)

		if err != nil {
			return rateLimitInfo, err
		}

		if !response.IsSuccess() {
			return rateLimitInfo, apperror.Internal(errorResponse.Detail)
		}

		if successResponse.Data.ProcessingInfo == nil {
			return rateLimitInfo, ErrFailedToFinalizeUpload
		}

		state := successResponse.Data.ProcessingInfo.State
		switch state {
		case helpers.ProcessingInfoStateSucceeded:
			return rateLimitInfo, nil
		case helpers.ProcessingInfoStateFailed:
			return rateLimitInfo, ErrFailedToFinalizeUpload
		case helpers.ProcessingInfoStateInProgress, helpers.ProcessingInfoStatePending:
			continue
		default:
			return rateLimitInfo, ErrFailedToFinalizeUpload
		}
	}
}

type TweetCreateResponse struct {
	Data struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"data"`
}

func publishTweetWithMedia(ctx context.Context, accessToken string, text string, mediaIDs []string, replyToTweetID string) (helpers.RateLimitInfo, error) {
	client := helpers.GetXClient(ctx, accessToken)

	tweetData := map[string]any{
		"text": text,
	}

	if len(mediaIDs) > 0 {
		tweetData["media"] = map[string]any{
			"media_ids": mediaIDs,
		}
	}

	if replyToTweetID != "" {
		tweetData["reply"] = map[string]any{
			"in_reply_to_tweet_id": replyToTweetID,
		}
	}

	var successResponse TweetCreateResponse
	var errorResponse helpers.ErrorResponse
	response, err := client.Post("/2/tweets").
		JSONContentType().
		Body(tweetData).
		Do(&successResponse, &errorResponse)

	rateLimitInfo := helpers.ExtractRateLimitInfo(response.Headers)

	if err != nil {
		return rateLimitInfo, err
	}

	if !response.IsSuccess() {
		return rateLimitInfo, apperror.Internal(errorResponse.Detail)
	}

	return rateLimitInfo, nil
}
