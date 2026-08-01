package token

import (
	authDomainToken "github.com/blocknextai/file-gateway-api/internal/auth/domain/token"
	"github.com/blocknextai/go-packages/auth/jwt"
	"github.com/blocknextai/go-packages/uuid"
)

type service struct {
	jwtService jwt.AuthJWTService
}

func NewService(jwtService jwt.AuthJWTService) authDomainToken.Service {
	return &service{jwtService: jwtService}
}

func (s *service) Generate() (string, error) {
	userID := uuid.NewV7()
	sessionID := uuid.NewV7()
	return s.jwtService.GenerateAccessToken(userID, sessionID)
}
