package credentialschema

import (
	"log/slog"
	"sync"

	"github.com/blocknextai/go-packages/json"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
)

type CredentialSchemaContextBuilder struct {
	credentialService nodeengineContract.CredentialService
	once              sync.Once
	cached            string
}

func NewCredentialSchemaContextBuilder(
	credentialService nodeengineContract.CredentialService,
) *CredentialSchemaContextBuilder {
	return &CredentialSchemaContextBuilder{
		credentialService: credentialService,
	}
}

func (b *CredentialSchemaContextBuilder) Build() string {
	b.once.Do(func() {
		credentials := b.credentialService.GetAllCredentials()

		data, err := json.Marshal(credentials)
		if err != nil {
			slog.Error("Failed to marshal credential schema context",
				"component", "Generation",
				"error", err,
			)
			b.cached = "[]"
			return
		}

		b.cached = string(data)
	})
	return b.cached
}
