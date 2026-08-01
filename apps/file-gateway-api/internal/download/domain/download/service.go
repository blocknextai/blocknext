package download

import (
	"context"
)

type Service interface {
	Download(ctx context.Context, rawURL string) (*FileInfo, error)
}
