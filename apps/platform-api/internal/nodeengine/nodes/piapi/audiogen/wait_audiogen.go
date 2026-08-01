package audiogen

import (
	"context"
	"time"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
)

const (
	WaitAudioGenMaxRetry = 30
)

var (
	WaitAudioGenErrorMap = map[string]string{
		"failed": "failed",
		"staged": "staged",
	}

	ErrWaitAudioGenMaxRetryReached = apperror.Internal("max retry reached")
)

type WaitAudioGenInput struct {
	TaskID string
}

type WaitAudioGenSuccessResponse struct {
	Code int `json:"code"`
	Data struct {
		Status string `json:"status"`
		Output struct {
			Songs []struct {
				ImagePath string `json:"image_path"`
				Lyrics    string `json:"lyrics"`
				SongPath  string `json:"song_path"`
			} `json:"songs"`
		} `json:"output"`
	} `json:"data"`
	Message string `json:"message"`
}

type WaitAudioGenErrorResponse struct {
	Message string `json:"message"`
}

func WaitAudioGen(ctx context.Context, client *httpclient.Client, input WaitAudioGenInput) ([]map[string]any, error) {
	retryDelayDuration := time.Duration(30_000) * time.Millisecond
	retryCount := 0

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(retryDelayDuration):
		}

		if retryCount++; retryCount > WaitAudioGenMaxRetry {
			return nil, ErrWaitAudioGenMaxRetryReached
		}

		var successResponse WaitAudioGenSuccessResponse
		var errorResponse WaitAudioGenErrorResponse

		response, err := client.Get("/task/"+input.TaskID).
			Do(&successResponse, &errorResponse)

		if err != nil {
			return nil, err
		}

		if !response.IsSuccess() {
			return nil, apperror.Internal(errorResponse.Message)
		}

		status := successResponse.Data.Status

		switch status {
		case "completed":
			items := make([]map[string]any, 0, len(successResponse.Data.Output.Songs))
			for _, songData := range successResponse.Data.Output.Songs {
				items = append(items, map[string]any{
					"audio": songData.SongPath,
					"image": songData.ImagePath,
					"text":  songData.Lyrics,
				})
			}
			return items, nil
		case "pending", "processing":
			continue
		default:
			if msg, ok := WaitAudioGenErrorMap[status]; ok {
				return nil, apperror.Internal(msg)
			}
			return nil, apperror.Internal("unknown status: " + status)
		}
	}
}
