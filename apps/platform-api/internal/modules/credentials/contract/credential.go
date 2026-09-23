package contract

import (
	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/modules/credentials/internal/application/credentials"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/modules/credentials/internal/domain/credentials"
)

type CredentialInfo = credentialsDomainCredentials.CredentialInfo
type SourceType = credentialsDomainCredentials.SourceType

const SourceTypePlatform = credentialsDomainCredentials.SourceTypePlatform

type CredentialService = credentialsApplicationCredentials.CredentialService

type Credential = credentialsDomainCredentials.Credential
