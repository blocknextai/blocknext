package imagegen

import (
	"context"
	"time"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
)

const (
	WaitImageGenMaxRetry = 30
)

var (
	WaitImageGenErrorMap = map[string]string{
		"failed": "failed",
		"staged": "staged",
	}

	ErrWaitImageGenMaxRetryReached = apperror.Internal("max retry reached")
)

type WaitImageGenInput struct {
	TaskID string
}

type WaitImageGenSuccessResponse struct {
	Code int `json:"code"`
	Data struct {
		Status string `json:"status"`
		Output struct {
			ImageURL    string `json:"image_url"`
			ImageBase64 string `json:"image_base64"`
		} `json:"output"`
	} `json:"data"`
	Message string `json:"message"`
}

type WaitImageGenErrorResponse struct {
	Message string `json:"message"`
}

func WaitImageGen(ctx context.Context, client *httpclient.Client, input WaitImageGenInput) (map[string]any, error) {
	retryDelayDuration := time.Duration(30_000) * time.Millisecond
	retryCount := 0

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(retryDelayDuration):
		}

		if retryCount++; retryCount > WaitImageGenMaxRetry {
			return nil, ErrWaitImageGenMaxRetryReached
		}

		var successResponse WaitImageGenSuccessResponse
		var errorResponse WaitImageGenErrorResponse

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
			return map[string]any{
				"image": successResponse.Data.Output.ImageURL,
			}, nil
		case "pending", "processing":
			continue
		default:
			if msg, ok := WaitImageGenErrorMap[status]; ok {
				return nil, apperror.Internal(msg)
			}
			return nil, apperror.Internal("unknown status: " + status)
		}
	}
}
