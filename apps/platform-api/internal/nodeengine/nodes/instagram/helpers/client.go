package helpers

import (
	"context"

	"github.com/blocknextai/go-packages/httpclient"
)

func GetInstagramClient(ctx context.Context, accessToken string) *httpclient.Client {
	return httpclient.NewClientBuilder().
		Context(ctx).
		BaseURL("https://graph.instagram.com/v23.0").
		QueryParam("access_token", accessToken).
		Build()
}
