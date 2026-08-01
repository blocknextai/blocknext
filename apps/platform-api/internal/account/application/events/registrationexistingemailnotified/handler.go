package registrationexistingemailnotified

import (
	"context"

	"github.com/blocknextai/platform-api/internal/account/application/auth/mailer"
	accountDomainEmails "github.com/blocknextai/platform-api/internal/account/domain/emails"
	"github.com/blocknextai/platform-api/internal/eventbus"
)

type Handler struct {
	mailer *mailer.Mailer
}

func New(mailer *mailer.Mailer, eventBus *eventbus.Bus) *Handler {
	handler := &Handler{mailer: mailer}
	eventbus.Subscribe(eventBus, handler.Handle)
	return handler
}

func (h *Handler) Handle(ctx context.Context, event accountDomainEmails.RegistrationExistingEmailNotifiedDomainEvent) error {
	return h.mailer.SendAlreadyRegistered(ctx, event.Email)
}
