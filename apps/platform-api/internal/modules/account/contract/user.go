package contract

import (
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/modules/account/internal/application/users"
	accountDomainUsers "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/users"
)

type EmailChangedDomainEvent = accountDomainUsers.EmailChangedDomainEvent
type PasswordChangedDomainEvent = accountDomainUsers.PasswordChangedDomainEvent
type User = accountDomainUsers.User
type UserCreatedDomainEvent = accountDomainUsers.UserCreatedDomainEvent

type UserService = accountApplicationUsers.UserService
