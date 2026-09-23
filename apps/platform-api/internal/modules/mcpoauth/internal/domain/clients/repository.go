package clients

import (
	"context"
)

type ClientRepository interface {
	GetByClientID(ctx context.Context, clientID string) (*Client, error)
	Create(ctx context.Context, client *Client) error
}
