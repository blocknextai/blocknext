package helpers

import (
	"context"

	"github.com/blocknextai/go-packages/httpclient"
)

func CreateClient(ctx context.Context, apiKey string) *httpclient.Client {
	return httpclient.NewClientBuilder().
		Context(ctx).
		BaseURL("https://api.anthropic.com/v1").
		Header("x-api-key", apiKey).
		Header("anthropic-version", "2023-06-01").
		Build()
}
