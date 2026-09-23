package eventbus

import (
	"context"
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/config"
	"github.com/blocknextai/platform-api/internal/eventbus/idempotency"
	eventbusPostgresInboxEntries "github.com/blocknextai/platform-api/internal/eventbus/postgres/inboxentries"
	eventbusPostgresOutboxMessages "github.com/blocknextai/platform-api/internal/eventbus/postgres/outboxmessages"
	"github.com/blocknextai/platform-api/internal/eventbus/publishing"
	"github.com/blocknextai/platform-api/internal/eventbus/relay"
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
	outboxRepository := eventbusPostgresOutboxMessages.NewRepository(deps.DB)
	inboxRepository := eventbusPostgresInboxEntries.NewRepository(deps.DB)

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
