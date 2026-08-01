package cqrs

import (
	"context"
)

type Handler[C any, R any] interface {
	Handle(ctx context.Context, command C) (R, error)
}

type HandlerFunc[C any, R any] func(ctx context.Context, command C) (R, error)

func (f HandlerFunc[C, R]) Handle(ctx context.Context, command C) (R, error) {
	return f(ctx, command)
}

type Behavior[C any, R any] func(next Handler[C, R]) Handler[C, R]

func Chain[C any, R any](handler Handler[C, R], behaviors ...Behavior[C, R]) Handler[C, R] {
	for i := len(behaviors) - 1; i >= 0; i-- {
		handler = behaviors[i](handler)
	}

	return handler
}
