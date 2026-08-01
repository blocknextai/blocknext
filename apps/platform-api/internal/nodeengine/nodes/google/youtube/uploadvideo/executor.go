package uploadvideo

import (
	"bytes"
	"context"
	"io"
	"strconv"
	"strings"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/youtube/helpers"
)

const (
	uploadChunkSize = 2 * 1024 * 1024
)

var (
	errFailedToReadVideoContent     = apperror.Internal("failed to read video content")
	errMissingUploadURL             = apperror.Internal("missing upload url")
	errUploadFailedChunk            = apperror.Internal("upload failed chunk")
	errUploadCompletedButNoResponse = apperror.Internal("upload completed but no response")
	errChunkUploadFailedWithDetails = apperror.Internal("chunk upload failed with details")
	errUploadFailedWithStatus       = apperror.Internal("upload failed with status")
)

type YouTubeUploadVideoExecutorInput struct {
	Title       string `schema:"title"`
	Description string `schema:"description"`
	CategoryID  string `schema:"categoryId"`
	Privacy     string `schema:"privacy"`
	VideoURL    string `schema:"videoUrl"`
}

type YouTubeUploadVideoExecutor struct {
	executors.Executor
	validator   *jsonschema.Validator[YouTubeUploadVideoExecutorInput]
	fileGateway filegateway.FileGateway
}

func NewYouTubeUploadVideoExecutor(
	nodeID string,
	validator *jsonschema.Validator[YouTubeUploadVideoExecutorInput],
	fileGateway filegateway.FileGateway,
) *YouTubeUploadVideoExecutor {
	return &YouTubeUploadVideoExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator:   validator,
		fileGateway: fileGateway,
	}
}

type YouTubeUploadSessionResponse struct {
	Kind string `json:"kind"`
	ETag string `json:"etag"`
}

type YouTubeVideoUploadResponse struct {
	Kind    string              `json:"kind"`
	ETag    string              `json:"etag"`
	ID      string              `json:"id"`
	Snippet YouTubeVideoSnippet `json:"snippet"`
	Status  YouTubeVideoStatus  `json:"status"`
	Details YouTubeVideoDetails `json:"contentDetails"`
}

type YouTubeVideoSnippet struct {
	PublishedAt          string                 `json:"publishedAt"`
	ChannelID            string                 `json:"channelId"`
	Title                string                 `json:"title"`
	Description          string                 `json:"description"`
	Thumbnails           YouTubeVideoThumbnails `json:"thumbnails"`
	ChannelTitle         string                 `json:"channelTitle"`
	CategoryID           string                 `json:"categoryId"`
	LiveBroadcastContent string                 `json:"liveBroadcastContent"`
	Localized            YouTubeVideoLocalized  `json:"localized"`
}

type YouTubeVideoThumbnails struct {
	Default  YouTubeVideoThumbnail `json:"default"`
	Medium   YouTubeVideoThumbnail `json:"medium"`
	High     YouTubeVideoThumbnail `json:"high"`
	Standard YouTubeVideoThumbnail `json:"standard"`
	Maxres   YouTubeVideoThumbnail `json:"maxres"`
}

type YouTubeVideoThumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type YouTubeVideoLocalized struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type YouTubeVideoStatus struct {
	UploadStatus        string `json:"uploadStatus"`
	PrivacyStatus       string `json:"privacyStatus"`
	License             string `json:"license"`
	Embeddable          bool   `json:"embeddable"`
	PublicStatsViewable bool   `json:"publicStatsViewable"`
}

type YouTubeVideoDetails struct {
	Duration        string `json:"duration"`
	Dimension       string `json:"dimension"`
	Definition      string `json:"definition"`
	Caption         string `json:"caption"`
	LicensedContent bool   `json:"licensedContent"`
	Projection      string `json:"projection"`
}

type YouTubeErrorResponse struct {
	Error YouTubeError `json:"error"`
}

type YouTubeError struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Errors  []YouTubeErrorDetail `json:"errors"`
	Status  string               `json:"status"`
}

type YouTubeErrorDetail struct {
	Message      string `json:"message"`
	Domain       string `json:"domain"`
	Reason       string `json:"reason"`
	Location     string `json:"location"`
	LocationType string `json:"locationType"`
}

func (e *YouTubeUploadVideoExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "youtube_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")
		client := helpers.CreateClient(ctx, accessToken)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			downloadResult, err := e.fileGateway.DownloadFile(ctx, input.VideoURL)
			if err != nil {
				return nil, err
			}
			defer func() {
				_ = downloadResult.DataReader.Close()
			}()

			videoData, err := io.ReadAll(downloadResult.DataReader)
			if err != nil {
				return nil, errFailedToReadVideoContent
			}

			contentType := downloadResult.ContentType
			if contentType == "" {
				contentType = "video/*"
			}

			totalSize := len(videoData)

			payload := map[string]any{
				"snippet": map[string]any{
					"title":       input.Title,
					"description": input.Description,
					"categoryId":  input.CategoryID,
				},
				"status": map[string]any{
					"privacyStatus": input.Privacy,
				},
			}

			var sessionResp YouTubeUploadSessionResponse
			var sessionErr YouTubeErrorResponse
			initReq := client.Post("/videos").
				JSONContentType().
				QueryParam("uploadType", "resumable").
				QueryParam("part", "snippet,status,contentDetails").
				Header("X-Upload-Content-Type", contentType).
				Header("X-Upload-Content-Length", strconv.Itoa(totalSize)).
				Body(payload)

			initRes, err := initReq.Do(&sessionResp, &sessionErr)
			if err != nil {
				return nil, err
			}

			if !initRes.IsSuccess() {
				return nil, apperror.Internal(sessionErr.Error.Message)
			}

			uploadURL := initRes.Headers.Get("Location")
			if uploadURL == "" {
				return nil, errMissingUploadURL
			}

			result, err := e.performChunkedUpload(ctx, uploadURL, videoData, totalSize, contentType, accessToken)
			if err != nil {
				return nil, err
			}

			results = append(results, map[string]any{
				"status":   true,
				"videoId":  result.ID,
				"videoUrl": "https://www.youtube.com/watch?v=" + result.ID,
			})
		}

		return results, nil
	}
}

func (e *YouTubeUploadVideoExecutor) performChunkedUpload(ctx context.Context, uploadURL string, videoData []byte, totalSize int, contentType, accessToken string) (*YouTubeVideoUploadResponse, error) {
	putClient := httpclient.NewClientBuilder().
		Context(ctx).
		BaseURL(uploadURL).
		BearerAuth(accessToken).
		Build()

	for start := 0; start < totalSize; start += uploadChunkSize {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		end := min(start+uploadChunkSize, totalSize)

		chunk := videoData[start:end]
		var rb strings.Builder
		rb.WriteString("bytes ")
		rb.WriteString(strconv.Itoa(start))
		rb.WriteByte('-')
		rb.WriteString(strconv.Itoa(end - 1))
		rb.WriteByte('/')
		rb.WriteString(strconv.Itoa(totalSize))
		rangeHeader := rb.String()

		var uploadResp YouTubeVideoUploadResponse
		var uploadErr YouTubeErrorResponse

		response, err := putClient.Put("").
			Header("Content-Length", strconv.Itoa(len(chunk))).
			Header("Content-Type", contentType).
			Header("Content-Range", rangeHeader).
			BodyReader(bytes.NewReader(chunk)).
			Do(&uploadResp, &uploadErr)

		if err != nil {
			return nil, errChunkUploadFailedWithDetails
		}

		switch response.Status {
		case 200, 201:
			if uploadResp.ID == "" {
				return nil, errUploadCompletedButNoResponse
			}
			return &uploadResp, nil
		case 308:
			continue
		case 416:
			rangeHeader := response.Headers.Get("Range")
			if rangeHeader != "" {
				if uploadedBytes := e.parseUploadedRange(rangeHeader); uploadedBytes > 0 {
					start = uploadedBytes - uploadChunkSize
				}
			}
			continue
		default:
			if uploadErr.Error.Code != 0 {
				return nil, errUploadFailedWithStatus
			}
			return nil, errUploadFailedChunk
		}
	}

	return nil, errUploadCompletedButNoResponse
}

func (e *YouTubeUploadVideoExecutor) parseUploadedRange(rangeHeader string) int {
	if after, ok := strings.CutPrefix(rangeHeader, "bytes="); ok {
		rangeStr := after
		parts := strings.Split(rangeStr, "-")
		if len(parts) == 2 {
			if end, err := strconv.Atoi(parts[1]); err == nil {
				return end + 1
			}
		}
	}
	return 0
}
