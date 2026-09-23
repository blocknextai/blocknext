package clients

import (
	"context"
)

type MetadataDocumentFetcher interface {
	Fetch(ctx context.Context, clientID string) (*Client, error)
}
