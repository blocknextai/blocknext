package createlyrics

import (
	"context"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
)

const (
	WaitLyricsMaxRetry = 30
)

var (
	WaitLyricsErrorMap = map[string]string{
		"CREATE_TASK_FAILED":    "create_task_failed",
		"GENERATE_LYRIC_FAILED": "generate_lyric_failed",
		"CALLBACK_EXCEPTION":    "callback_exception",
		"SENSITIVE_WORD_ERROR":  "sensitive_word_error",
	}

	ErrWaitLyricsMaxRetryReached = apperror.Internal("max retry reached")
)

type WaitLyricsInput struct {
	TaskID string
}

type WaitLyricsSuccessResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		TaskID   string `json:"taskId"`
		Param    string `json:"param"`
		Response struct {
			TaskID string `json:"taskId"`
			Data   []struct {
				Text         string `json:"text"`
				Title        string `json:"title"`
				Status       string `json:"status"`
				ErrorMessage string `json:"errorMessage"`
			} `json:"data"`
		} `json:"response"`
		Status       string `json:"status"`
		Type         string `json:"type"`
		ErrorCode    any    `json:"errorCode"`
		ErrorMessage any    `json:"errorMessage"`
	} `json:"data"`
}

type WaitLyricsErrorResponse struct{}

func WaitLyrics(ctx context.Context, client *httpclient.Client, input WaitLyricsInput) (map[string]any, error) {
	retryDelayDuration := time.Duration(20_000) * time.Millisecond
	retryCount := 0

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(retryDelayDuration):
		}

		if retryCount++; retryCount > WaitLyricsMaxRetry {
			return nil, ErrWaitLyricsMaxRetryReached
		}

		var successResponse WaitLyricsSuccessResponse
		var errorResponse WaitLyricsErrorResponse

		_, err := client.Get("/lyrics/record-info").
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
			contents := []string{}
			for _, data := range successResponse.Data.Response.Data {
				contents = append(contents, data.Text)
			}
			content := strings.Join(contents, "\n\n")

			return map[string]any{
				"text": content,
			}, nil
		case "PENDING":
			continue
		default:
			if msg, ok := WaitLyricsErrorMap[status]; ok {
				return nil, apperror.Internal(msg)
			}
			return nil, apperror.Internal("unknown status: " + status)
		}
	}
}
