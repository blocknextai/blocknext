package userwelcome

import (
	"context"

	"github.com/blocknextai/platform-api/internal/account/application/auth/mailer"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/account/domain/users"
	"github.com/blocknextai/platform-api/internal/eventbus"
)

type Handler struct {
	mailer *mailer.Mailer
}

func New(mailer *mailer.Mailer, eventBus *eventbus.Bus) *Handler {
	handler := &Handler{
		mailer: mailer,
	}
	eventbus.Subscribe(eventBus, handler.Handle)
	return handler
}

func (h *Handler) Handle(ctx context.Context, event accountDomainUsers.UserCreatedDomainEvent) error {
	return h.mailer.SendWelcome(ctx, event.Identifier)
}
