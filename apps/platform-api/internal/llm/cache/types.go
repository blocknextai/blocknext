package cache

import (
	"context"
)

type Cache interface {
	Ensure(ctx context.Context, instruction string) string
}
