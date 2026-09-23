package helpers

import (
	"context"

	"github.com/blocknextai/go-packages/httpclient"
)

func CreateClient(ctx context.Context, apiKey string) *httpclient.Client {
	return httpclient.NewClientBuilder().
		Context(ctx).
		BaseURL("https://api.coingecko.com/api/v3").
		Header("X-CG-Pro-API-Key", apiKey).
		Build()
}
