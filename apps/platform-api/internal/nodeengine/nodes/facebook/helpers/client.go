package helpers

import (
	"context"

	"github.com/blocknextai/go-packages/httpclient"
)

func GetFacebookClient(ctx context.Context, accessToken string) *httpclient.Client {
	return httpclient.NewClientBuilder().
		Context(ctx).
		BaseURL("https://graph.facebook.com/v23.0").
		QueryParam("access_token", accessToken).
		Build()
}
