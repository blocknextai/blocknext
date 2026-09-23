package helpers

import (
	"context"

	"github.com/blocknextai/go-packages/httpclient"
)

func GetXClient(ctx context.Context, accessToken string) *httpclient.Client {
	return httpclient.NewClientBuilder().
		Context(ctx).
		BaseURL("https://api.x.com").
		BearerAuth(accessToken).
		Build()
}
