package token

import (
	"context"
	"time"
)

const (
	BaseDelay = 50 * time.Millisecond
	MaxDelay  = 2 * time.Second
)

func AcquireWithBackoff(ctx context.Context, tryAcquire func(ctx context.Context) (bool, error)) (chan struct{}, error) {
	currentDelay := BaseDelay

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		acquired, err := tryAcquire(ctx)
		if err != nil {
			return nil, err
		}

		if acquired {
			token := make(chan struct{}, 1)
			token <- struct{}{}
			return token, nil
		}

		select {
		case <-time.After(currentDelay):
			if currentDelay < MaxDelay {
				currentDelay = min(currentDelay*2, MaxDelay)
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func Drain(token chan struct{}) {
	select {
	case <-token:
	default:
	}
}
