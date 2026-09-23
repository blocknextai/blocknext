package helpers

import (
	"context"

	"github.com/blocknextai/go-packages/httpclient"
)

func CreateClient(ctx context.Context, apiKey string) *httpclient.Client {
	baseURL := "https://api-free.deepl.com/v2"
	if len(apiKey) > 0 && apiKey[len(apiKey)-3:] != ":fx" {
		baseURL = "https://api.deepl.com/v2"
	}

	return httpclient.NewClientBuilder().
		Context(ctx).
		BaseURL(baseURL).
		Header("Authorization", "DeepL-Auth-Key "+apiKey).
		Build()
}
