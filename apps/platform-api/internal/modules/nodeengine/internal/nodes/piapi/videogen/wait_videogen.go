package videogen

import (
	"context"
	"time"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
)

const (
	WaitVideoGenMaxRetry = 30
)

var (
	WaitVideoGenErrorMap = map[string]string{
		"failed": "failed",
		"staged": "staged",
	}

	ErrWaitVideoGenMaxRetryReached = apperror.Internal("max retry reached")
)

type WaitVideoGenInput struct {
	TaskID string
}

type WaitVideoGenSuccessResponse struct {
	Code int `json:"code"`
	Data struct {
		Status string `json:"status"`
		Output struct {
			Works []struct {
				Cover struct {
					Resource                 string `json:"resource"`
					ResourceWithoutWatermark string `json:"resource_without_watermark"`
				} `json:"cover"`
				Video struct {
					Resource                 string `json:"resource"`
					ResourceWithoutWatermark string `json:"resource_without_watermark"`
				} `json:"video"`
			} `json:"works"`
		} `json:"output"`
	} `json:"data"`
	Message string `json:"message"`
}

type WaitVideoGenErrorResponse struct {
	Message string `json:"message"`
}

func WaitVideoGen(ctx context.Context, client *httpclient.Client, input WaitVideoGenInput) ([]map[string]any, error) {
	retryDelayDuration := time.Duration(30_000) * time.Millisecond
	retryCount := 0

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(retryDelayDuration):
		}

		if retryCount++; retryCount > WaitVideoGenMaxRetry {
			return nil, ErrWaitVideoGenMaxRetryReached
		}

		var successResponse WaitVideoGenSuccessResponse
		var errorResponse WaitVideoGenErrorResponse

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
			items := make([]map[string]any, 0, len(successResponse.Data.Output.Works))
			for _, workData := range successResponse.Data.Output.Works {
				items = append(items, map[string]any{
					"image": workData.Cover.Resource,
					"video": workData.Video.Resource,
				})
			}
			return items, nil
		case "pending", "processing":
			continue
		default:
			if msg, ok := WaitVideoGenErrorMap[status]; ok {
				return nil, apperror.Internal(msg)
			}
			return nil, apperror.Internal("unknown status: " + status)
		}
	}
}
