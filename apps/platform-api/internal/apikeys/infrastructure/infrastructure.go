package infrastructure

import (
	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/apikeys/application/apikeys/createapikey"
	"github.com/blocknextai/platform-api/internal/apikeys/application/apikeys/deleteapikey"
	"github.com/blocknextai/platform-api/internal/apikeys/application/apikeys/getallapikeys"
	"github.com/blocknextai/platform-api/internal/apikeys/application/apikeys/getscopes"
	"github.com/blocknextai/platform-api/internal/apikeys/application/apikeys/regenerateapikey"
	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/apikeys/domain/apikeys"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
)

type Handlers struct {
	CreateAPIKey     cqrs.Handler[*createapikey.CreateAPIKeyCommand, *createapikey.CreateAPIKeyResponse]
	GetAllAPIKeys    cqrs.Handler[*getallapikeys.GetAllAPIKeysQuery, *getallapikeys.GetAllAPIKeysResponse]
	GetScopes        cqrs.Handler[*getscopes.GetScopesQuery, *getscopes.GetScopesResponse]
	DeleteAPIKey     cqrs.Handler[*deleteapikey.DeleteAPIKeyCommand, *deleteapikey.DeleteAPIKeyResponse]
	RegenerateAPIKey cqrs.Handler[*regenerateapikey.RegenerateAPIKeyCommand, *regenerateapikey.RegenerateAPIKeyResponse]
}

type RegisterInfrastructureDeps struct {
	TransactionManager database.TransactionManager

	ApiKeyRepository apiKeysDomainAPIKeys.APIKeyRepository
}

func RegisterInfrastructure(deps RegisterInfrastructureDeps) *Handlers {
	return &Handlers{
		CreateAPIKey:     cqrs.ValidationBehavior(createapikey.New(deps.ApiKeyRepository, deps.TransactionManager)),
		GetAllAPIKeys:    cqrs.ValidationBehavior(getallapikeys.New(deps.ApiKeyRepository)),
		GetScopes:        cqrs.ValidationBehavior(getscopes.New()),
		DeleteAPIKey:     cqrs.ValidationBehavior(deleteapikey.New(deps.ApiKeyRepository, deps.TransactionManager)),
		RegenerateAPIKey: cqrs.ValidationBehavior(regenerateapikey.New(deps.ApiKeyRepository, deps.TransactionManager)),
	}
}
