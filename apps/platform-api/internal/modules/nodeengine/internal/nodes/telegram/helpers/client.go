package helpers

import (
	"context"

	"github.com/blocknextai/go-packages/httpclient"
)

func CreateClient(ctx context.Context, botToken string) *httpclient.Client {
	return httpclient.NewClientBuilder().
		Context(ctx).
		BaseURL("https://api.telegram.org/bot" + botToken).
		Build()
}
