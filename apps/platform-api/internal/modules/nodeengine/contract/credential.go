package contract

import (
	nodeengineApplicationCredentials "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/credentials"
	nodeengineDomainCredentials "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/credentials"
)

type CredentialIcon = nodeengineDomainCredentials.CredentialIcon
type CredentialManager = nodeengineDomainCredentials.CredentialManager

type CredentialProcessor = nodeengineApplicationCredentials.CredentialProcessor
type CredentialService = nodeengineApplicationCredentials.CredentialService
