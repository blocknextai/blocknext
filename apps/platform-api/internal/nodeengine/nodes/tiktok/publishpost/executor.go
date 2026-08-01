package publishpost

import (
	"bytes"
	"context"
	"io"
	"strconv"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/tiktok/helpers"
)

const (
	minChunkSize      = 5 * 1024 * 1024
	maxChunkSize      = 64 * 1024 * 1024
	maxFinalChunkSize = 128 * 1024 * 1024
)

var (
	ErrVideoTooLarge        = apperror.Internal("video too large")
	ErrUploadURLNotProvided = apperror.Internal("upload url not provided")
	ErrChunkSizeTooSmall    = apperror.Internal("chunk size too small")
	ErrChunkSizeTooLarge    = apperror.Internal("chunk size too large")
	ErrFinalChunkTooLarge   = apperror.Internal("final chunk too large")
	ErrUploadFailed         = apperror.Internal("upload failed")
)

type TiktokPublishPostExecutorInput struct {
	VideoURL string `schema:"videoUrl"`
	Caption  string `schema:"caption"`
}

type TiktokPublishPostExecutor struct {
	executors.Executor
	validator   *jsonschema.Validator[TiktokPublishPostExecutorInput]
	fileGateway filegateway.FileGateway
}

func NewTiktokPublishPostExecutor(
	nodeID string,
	validator *jsonschema.Validator[TiktokPublishPostExecutorInput],
	fileGateway filegateway.FileGateway,
) *TiktokPublishPostExecutor {
	return &TiktokPublishPostExecutor{
		Executor: executors.Executor{
			ID: nodeID,
		},
		validator:   validator,
		fileGateway: fileGateway,
	}
}

type UploadInitResponse struct {
	Data struct {
		PublishID string `json:"publish_id"`
		UploadURL string `json:"upload_url"`
	} `json:"data"`
	Error ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (e *TiktokPublishPostExecutor) ExecuteWithContext(ctx context.Context, credentials map[string]any, data []map[string]any) ([]map[string]any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		credential := nodeEngineDomainCredentials.GetCredentials(credentials, "tiktok_oauth2")
		oauthToken := credential.Object("oauthToken")
		accessToken := oauthToken.String("accessToken")
		client := helpers.CreateClient(ctx, accessToken)

		results := make([]map[string]any, 0)
		for _, item := range data {
			input, err := e.validator.Parse(item)
			if err != nil {
				return nil, err
			}

			publishID, uploadError := e.uploadVideo(ctx, client, input.VideoURL, input.Caption)
			if uploadError != nil {
				return nil, uploadError
			}

			results = append(results, map[string]any{
				"status":    true,
				"publishId": publishID,
			})
		}

		return results, nil
	}
}

func (e *TiktokPublishPostExecutor) uploadVideo(
	ctx context.Context,
	client *httpclient.Client,
	videoURL string,
	caption string,
) (string, error) {
	downloadResult, err := e.fileGateway.DownloadFile(ctx, videoURL)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = downloadResult.DataReader.Close()
	}()

	videoData, err := io.ReadAll(downloadResult.DataReader)
	if err != nil {
		return "", err
	}

	contentType := downloadResult.ContentType
	if contentType == "" {
		contentType = "video/mp4"
	}

	postInfo := map[string]any{
		"privacy_level":            "PUBLIC_TO_EVERYONE",
		"video_cover_timestamp_ms": 1000,
	}

	if caption != "" {
		postInfo["title"] = caption
	}

	videoSize := int64(len(videoData))
	var chunkSize int64
	var totalChunks int64

	if videoSize < minChunkSize {
		chunkSize = videoSize
		totalChunks = 1
	} else {
		chunkSize = maxChunkSize
		totalChunks = videoSize / chunkSize
		if videoSize%chunkSize != 0 {
			totalChunks++
		}
		if totalChunks > 1000 {
			return "", ErrVideoTooLarge
		}
	}

	requestBody := map[string]any{
		"post_info": postInfo,
		"source_info": map[string]any{
			"source":            "FILE_UPLOAD",
			"video_size":        videoSize,
			"chunk_size":        chunkSize,
			"total_chunk_count": totalChunks,
		},
	}

	var uploadSuccessResponse UploadInitResponse
	var uploadErrorResponse ErrorResponse
	uploadResponse, err := client.Post("/post/publish/inbox/video/init/").
		JSONContentType().
		Body(requestBody).
		Do(&uploadSuccessResponse, &uploadErrorResponse)

	if err != nil {
		return "", err
	}

	if !uploadResponse.IsSuccess() {
		return "", apperror.Internal(uploadErrorResponse.Error.Message)
	}

	if uploadSuccessResponse.Data.UploadURL == "" {
		return "", ErrUploadURLNotProvided
	}

	err = uploadVideoChunks(uploadSuccessResponse.Data.UploadURL, videoData, chunkSize, contentType)
	if err != nil {
		return "", err
	}

	return uploadSuccessResponse.Data.PublishID, nil
}

func uploadVideoChunks(uploadURL string, videoData []byte, chunkSize int64, contentType string) error {
	totalSize := int64(len(videoData))
	client := helpers.CreateUploadClient()

	chunkCount := int64(0)
	for start := int64(0); start < totalSize; start += chunkSize {
		chunkCount++
		end := min(start+chunkSize, totalSize)

		isLastChunk := end == totalSize
		currentChunkSize := end - start

		if !isLastChunk && currentChunkSize < minChunkSize {
			return ErrChunkSizeTooSmall
		}

		if !isLastChunk && currentChunkSize > maxChunkSize {
			return ErrChunkSizeTooLarge
		}

		if isLastChunk && currentChunkSize > maxFinalChunkSize {
			return ErrFinalChunkTooLarge
		}

		chunkData := videoData[start:end]

		resp, err := client.Put(uploadURL).
			BodyReader(bytes.NewReader(chunkData)).
			Header("Content-Type", contentType).
			Header("Content-Length", strconv.Itoa(len(chunkData))).
			Header("Content-Range", "bytes "+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end-1, 10)+"/"+strconv.FormatInt(totalSize, 10)).
			DoRaw()
		if err != nil {
			return err
		}

		expectedStatus := 206
		if isLastChunk {
			expectedStatus = 201
		}

		if resp.Status != expectedStatus {
			return ErrUploadFailed
		}
	}

	return nil
}
