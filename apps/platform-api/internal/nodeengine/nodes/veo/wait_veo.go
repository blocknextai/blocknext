package veo

import (
	"context"
	"time"

	"github.com/blocknextai/go-packages/apperror"

	"github.com/blocknextai/go-packages/httpclient"
)

const (
	WaitVeoMaxRetry = 60
)

var (
	ErrWaitVeoMaxRetryReached   = apperror.Internal("max retry reached")
	ErrWaitVeoApiError          = apperror.Internal("api error")
	ErrWaitVeoStatusCheckFailed = apperror.Internal("status check failed")
	ErrWaitVeoNoVideoSamples    = apperror.Internal("no video samples")
)

type WaitVeoInput struct {
	OperationName string
	APIKey        string
}

type WaitVeoSuccessResponse struct {
	Name     string `json:"name"`
	Done     bool   `json:"done"`
	Response struct {
		GenerateVideoResponse struct {
			GeneratedSamples []struct {
				Video struct {
					URI string `json:"uri"`
				} `json:"video"`
			} `json:"generatedSamples"`
		} `json:"generateVideoResponse"`
	} `json:"response"`
}

type WaitVeoErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func WaitVeo(ctx context.Context, client *httpclient.Client, input WaitVeoInput) ([]string, error) {
	retryDelayDuration := time.Duration(10_000) * time.Millisecond
	retryCount := 0

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(retryDelayDuration):
		}

		if retryCount++; retryCount > WaitVeoMaxRetry {
			return nil, ErrWaitVeoMaxRetryReached
		}

		var successResponse WaitVeoSuccessResponse
		var errorResponse WaitVeoErrorResponse

		response, err := client.Get("/"+input.OperationName).
			QueryParam("key", input.APIKey).
			Do(&successResponse, &errorResponse)

		if err != nil {
			return nil, err
		}

		if !response.IsSuccess() {
			if errorResponse.Error.Message != "" {
				return nil, ErrWaitVeoApiError
			}
			return nil, ErrWaitVeoStatusCheckFailed
		}

		if successResponse.Done {
			samples := successResponse.Response.GenerateVideoResponse.GeneratedSamples
			if len(samples) > 0 {
				videoURLs := make([]string, 0, len(samples))
				for _, sample := range samples {
					videoURLs = append(videoURLs, sample.Video.URI+"&key="+input.APIKey)
				}
				return videoURLs, nil
			}
			return nil, ErrWaitVeoNoVideoSamples
		}

		continue
	}
}
