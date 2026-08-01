package eventbus

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/config"
	"github.com/blocknextai/platform-api/internal/eventbus/application/idempotency"
	"github.com/blocknextai/platform-api/internal/eventbus/application/publishing"
	"github.com/blocknextai/platform-api/internal/eventbus/application/relay"
	inboxEntriesInfrastructure "github.com/blocknextai/platform-api/internal/eventbus/infrastructure/inboxentries"
	outboxMessagesInfrastructure "github.com/blocknextai/platform-api/internal/eventbus/infrastructure/outboxmessages"
)

type Dependencies struct {
	DB                 *sql.DB
	TransactionManager database.TransactionManager

	Options config.EventBusOptions
}

type Module struct {
	PublisherService publishing.PublisherService
	InboxService     *idempotency.InboxService
	Bus              *Bus

	relay *relay.RelayService
}

func NewModule(deps Dependencies) *Module {
	outboxRepository := outboxMessagesInfrastructure.NewRepository(deps.DB)
	inboxRepository := inboxEntriesInfrastructure.NewRepository(deps.DB)

	bus := NewBus()

	relayService := relay.NewRelayService(outboxRepository, bus, relay.Options{
		PollInterval:    deps.Options.PollInterval,
		BatchSize:       deps.Options.BatchSize,
		MaxAttempts:     deps.Options.MaxAttempts,
		BaseBackoff:     deps.Options.BaseBackoff,
		StuckTimeout:    deps.Options.StuckTimeout,
		ReclaimInterval: deps.Options.ReclaimInterval,
	})

	return &Module{
		PublisherService: publishing.NewPublisherService(outboxRepository),
		InboxService:     idempotency.NewInboxService(inboxRepository, deps.TransactionManager),
		Bus:              bus,
		relay:            relayService,
	}
}

func (m *Module) StartRelay(ctx context.Context) {
	go m.relay.Start(ctx)
}
