package contract

import (
	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/application/linkedaccounts"
	accountDomainLinkedAccounts "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/linkedaccounts"
)

type LinkedAccount = accountDomainLinkedAccounts.LinkedAccount

type LinkedAccountService = accountApplicationLinkedAccounts.LinkedAccountService
