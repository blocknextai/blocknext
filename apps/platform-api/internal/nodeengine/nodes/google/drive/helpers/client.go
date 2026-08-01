package helpers

import (
	"context"

	"github.com/blocknextai/go-packages/httpclient"
)

func CreateClient(ctx context.Context, accessToken string) *httpclient.Client {
	return httpclient.NewClientBuilder().
		Context(ctx).
		BearerAuth(accessToken).
		BaseURL("https://www.googleapis.com/drive/v3").
		Build()
}
