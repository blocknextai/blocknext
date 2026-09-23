package token

import (
	"github.com/blocknextai/go-packages/auth/jwt"
	"github.com/blocknextai/go-packages/uuid"
)

type Service interface {
	Generate() (string, error)
}

type service struct {
	jwtService jwt.AuthJWTService
}

func NewService(jwtService jwt.AuthJWTService) Service {
	return &service{jwtService: jwtService}
}

func (s *service) Generate() (string, error) {
	userID := uuid.NewV7()
	sessionID := uuid.NewV7()
	return s.jwtService.GenerateAccessToken(userID, sessionID)
}
