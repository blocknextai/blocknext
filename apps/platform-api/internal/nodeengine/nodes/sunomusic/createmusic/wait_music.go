package createmusic

import (
	"context"
	"time"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
)

const (
	WaitMusicMaxRetry = 30
)

var (
	WaitMusicErrorMap = map[string]string{
		"CREATE_TASK_FAILED":    "create_task_failed",
		"GENERATE_AUDIO_FAILED": "generate_audio_failed",
		"CALLBACK_EXCEPTION":    "callback_exception",
		"SENSITIVE_WORD_ERROR":  "sensitive_word_error",
	}

	ErrWaitMusicMaxRetryReached = apperror.Internal("max retry reached")
)

type WaitMusicInput struct {
	TaskID string
}

type WaitMusicSuccessResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		TaskID        string `json:"taskId"`
		ParentMusicID string `json:"parentMusicId"`
		Param         string `json:"param"`
		Response      struct {
			TaskID   string `json:"taskId"`
			SunoData []struct {
				ID             string `json:"id"`
				AudioURL       string `json:"audioUrl"`
				StreamAudioURL string `json:"streamAudioUrl"`
				ImageURL       string `json:"imageUrl"`
				Prompt         string `json:"prompt"`
				ModelName      string `json:"modelName"`
				Title          string `json:"title"`
			} `json:"sunoData"`
		} `json:"response"`
		Status       string `json:"status"`
		Type         string `json:"type"`
		ErrorCode    any    `json:"errorCode"`
		ErrorMessage any    `json:"errorMessage"`
	} `json:"data"`
}

type WaitMusicErrorResponse struct{}

func WaitMusic(ctx context.Context, client *httpclient.Client, input WaitMusicInput) ([]map[string]any, error) {
	retryDelayDuration := time.Duration(20_000) * time.Millisecond
	retryCount := 0

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(retryDelayDuration):
		}

		if retryCount++; retryCount > WaitMusicMaxRetry {
			return nil, ErrWaitMusicMaxRetryReached
		}

		var successResponse WaitMusicSuccessResponse
		var errorResponse WaitMusicErrorResponse

		_, err := client.Get("/generate/record-info").
			QueryParam("taskId", input.TaskID).
			Do(&successResponse, &errorResponse)

		if err != nil {
			return nil, err
		}

		if successResponse.Code != 200 {
			return nil, apperror.Internal(successResponse.Msg)
		}

		status := successResponse.Data.Status
		switch status {
		case "SUCCESS":
			items := make([]map[string]any, 0, len(successResponse.Data.Response.SunoData))
			for _, data := range successResponse.Data.Response.SunoData {
				items = append(items, map[string]any{
					"audio": data.AudioURL,
					"image": data.ImageURL,
				})
			}
			return items, nil
		case "PENDING", "TEXT_SUCCESS", "FIRST_SUCCESS":
			continue
		default:
			if msg, ok := WaitMusicErrorMap[status]; ok {
				return nil, apperror.Internal(msg)
			}
			return nil, apperror.Internal("unknown status: " + status)
		}
	}
}
