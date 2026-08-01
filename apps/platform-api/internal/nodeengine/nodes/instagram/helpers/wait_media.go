package helpers

import (
	"context"
	"time"

	"github.com/blocknextai/go-packages/apperror"
)

const (
	WaitMediaMaxRetry = 30
)

var (
	ErrWaitMediaMaxRetryReached    = apperror.Internal("max retry reached")
	ErrWaitMediaUnhandledException = apperror.Internal("unhandled exception")
)

type MediaStatusResponse struct {
	StatusCode   string `json:"status_code"`
	ErrorMessage string `json:"error_message"`
}

type MediaStatusErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func WaitForMediaReady(ctx context.Context, accessToken string, mediaID string) error {
	retryDelayDuration := time.Duration(10_000) * time.Millisecond
	retryCount := 0

	client := GetInstagramClient(ctx, accessToken)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(retryDelayDuration):
		}

		if retryCount++; retryCount > WaitMediaMaxRetry {
			return ErrWaitMediaMaxRetryReached
		}

		var statusResponse MediaStatusResponse
		var errorResponse MediaStatusErrorResponse
		response, err := client.Get("/"+mediaID).
			Do(&statusResponse, &errorResponse)

		if err != nil {
			return err
		}

		if !response.IsSuccess() {
			return apperror.Internal(errorResponse.Error.Message)
		}

		switch statusResponse.StatusCode {
		case "FINISHED", "PUBLISHED":
			return nil
		case "ERROR", "EXPIRED":
			return apperror.Internal(statusResponse.ErrorMessage)
		case "IN_PROGRESS":
			continue
		default:
			return ErrWaitMediaUnhandledException
		}
	}
}
