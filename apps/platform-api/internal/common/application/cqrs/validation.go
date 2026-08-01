package cqrs

import (
	"context"
)

type Validatable interface {
	Validate() error
}

func ValidationBehavior[C any, R any](next Handler[C, R]) Handler[C, R] {
	return HandlerFunc[C, R](func(ctx context.Context, command C) (R, error) {
		if validatable, ok := any(command).(Validatable); ok {
			if err := validatable.Validate(); err != nil {
				var zero R
				return zero, err
			}
		}

		return next.Handle(ctx, command)
	})
}
