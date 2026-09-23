package magiclinkcreated

import (
	"context"

	"github.com/blocknextai/platform-api/internal/eventbus"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/application/auth/mailer"
	accountDomainEmails "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/emails"
)

type Handler struct {
	mailer *mailer.Mailer
}

func New(mailer *mailer.Mailer, eventBus *eventbus.Bus) *Handler {
	handler := &Handler{mailer: mailer}
	eventbus.Subscribe(eventBus, handler.Handle)
	return handler
}

func (h *Handler) Handle(ctx context.Context, event accountDomainEmails.MagicLinkCreatedDomainEvent) error {
	return h.mailer.SendMagicLink(ctx, event.Email, event.Token)
}
