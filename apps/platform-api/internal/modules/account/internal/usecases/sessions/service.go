package sessions

import (
	accountApplicationSessions "github.com/blocknextai/platform-api/internal/modules/account/internal/application/sessions"
	accountDomainSessions "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/sessions"
)

type Service struct {
	sessionRepository accountDomainSessions.SessionRepository
	sessionService    accountApplicationSessions.SessionService
}

func NewService(
	sessionRepository accountDomainSessions.SessionRepository,
	sessionService accountApplicationSessions.SessionService,
) *Service {
	return &Service{
		sessionRepository: sessionRepository,
		sessionService:    sessionService,
	}
}
