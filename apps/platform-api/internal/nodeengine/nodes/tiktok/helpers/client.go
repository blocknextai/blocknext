package helpers

import (
	"context"

	"github.com/blocknextai/go-packages/httpclient"
)

func CreateClient(ctx context.Context, accessToken string) *httpclient.Client {
	return httpclient.NewClientBuilder().
		Context(ctx).
		BaseURL("https://open.tiktokapis.com/v2").
		BearerAuth(accessToken).
		Build()
}

func CreateUploadClient() *httpclient.Client {
	return httpclient.NewClientBuilder().Build()
}
