package helpers

import (
	"context"

	"github.com/blocknextai/go-packages/httpclient"
)

func CreateClient(ctx context.Context, accessToken string) *httpclient.Client {
	return httpclient.NewClientBuilder().
		Context(ctx).
		BaseURL("https://api.notion.com/v1").
		BearerAuth(accessToken).
		Header("Notion-Version", "2022-06-28").
		Build()
}
