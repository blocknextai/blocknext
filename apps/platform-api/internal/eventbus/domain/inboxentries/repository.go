package inboxentries

import (
	"context"
)

type Repository interface {
	MarkProcessed(ctx context.Context, entry *InboxEntry) (bool, error)
}
