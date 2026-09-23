package helpers

import (
	"context"

	"github.com/blocknextai/go-packages/httpclient"
)

func CreateClient(ctx context.Context, apiKey string) *httpclient.Client {
	return httpclient.NewClientBuilder().
		Context(ctx).
		BaseURL("https://api.openai.com/v1").
		BearerAuth(apiKey).
		Build()
}
