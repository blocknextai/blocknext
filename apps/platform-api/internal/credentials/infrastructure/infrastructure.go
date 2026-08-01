package infrastructure

import (
	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/secretmanager"
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	"github.com/blocknextai/platform-api/internal/credentials/application/credentials/createcredential"
	"github.com/blocknextai/platform-api/internal/credentials/application/credentials/deletecredential"
	"github.com/blocknextai/platform-api/internal/credentials/application/credentials/getallcredentials"
	"github.com/blocknextai/platform-api/internal/credentials/application/credentials/getcredentialbyid"
	"github.com/blocknextai/platform-api/internal/credentials/application/credentials/getcredentialsfornodes"
	"github.com/blocknextai/platform-api/internal/credentials/application/credentials/updatecredential"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/credentials/domain/credentials"
	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
	platformApplicationPlatformCredentials "github.com/blocknextai/platform-api/internal/platform/application/platformcredentials"
)

type Handlers struct {
	CreateCredential       cqrs.Handler[*createcredential.CreateCredentialCommand, *createcredential.CreateCredentialResponse]
	GetAllCredentials      cqrs.Handler[*getallcredentials.GetAllCredentialsQuery, *getallcredentials.GetAllCredentialsResponse]
	GetCredentialByID      cqrs.Handler[*getcredentialbyid.GetCredentialByIDQuery, *getcredentialbyid.GetCredentialByIDResponse]
	GetCredentialsForNodes cqrs.Handler[*getcredentialsfornodes.GetCredentialsForNodesQuery, *getcredentialsfornodes.GetCredentialsForNodesResponse]
	UpdateCredential       cqrs.Handler[*updatecredential.UpdateCredentialCommand, *updatecredential.UpdateCredentialResponse]
	DeleteCredential       cqrs.Handler[*deletecredential.DeleteCredentialCommand, *deletecredential.DeleteCredentialResponse]
}

type RegisterInfrastructureDeps struct {
	TransactionManager database.TransactionManager
	SecretManager      secretmanager.SecretManager

	CredentialRepository        credentialsDomainCredentials.CredentialRepository
	CredentialProcessor         nodeEngineApplicationCredentials.CredentialProcessor
	NodeEngineCredentialService nodeEngineApplicationCredentials.CredentialService
	PlatformCredentialService   platformApplicationPlatformCredentials.PlatformCredentialService
}

func RegisterInfrastructure(deps RegisterInfrastructureDeps) *Handlers {
	return &Handlers{
		CreateCredential:       cqrs.ValidationBehavior(createcredential.New(deps.CredentialRepository, deps.SecretManager, deps.TransactionManager, deps.PlatformCredentialService)),
		GetAllCredentials:      cqrs.ValidationBehavior(getallcredentials.New(deps.CredentialRepository)),
		GetCredentialByID:      cqrs.ValidationBehavior(getcredentialbyid.New(deps.CredentialRepository, deps.SecretManager, deps.CredentialProcessor, deps.NodeEngineCredentialService)),
		GetCredentialsForNodes: cqrs.ValidationBehavior(getcredentialsfornodes.New(deps.CredentialRepository, deps.NodeEngineCredentialService)),
		UpdateCredential:       cqrs.ValidationBehavior(updatecredential.New(deps.CredentialRepository, deps.SecretManager, deps.TransactionManager, deps.CredentialProcessor, deps.NodeEngineCredentialService)),
		DeleteCredential:       cqrs.ValidationBehavior(deletecredential.New(deps.CredentialRepository, deps.TransactionManager)),
	}
}
