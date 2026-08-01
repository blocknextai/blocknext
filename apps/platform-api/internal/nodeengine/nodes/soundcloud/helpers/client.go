package helpers

import (
	"context"

	"github.com/blocknextai/go-packages/httpclient"
)

func CreateClient(ctx context.Context, accessToken string) *httpclient.Client {
	return httpclient.NewClientBuilder().
		Context(ctx).
		BaseURL("https://api.soundcloud.com").
		Header("Authorization", "OAuth "+accessToken).
		Build()
}
