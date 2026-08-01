package oauth2

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/blocknextai/go-packages/base64"
	"github.com/blocknextai/go-packages/cache"
	"github.com/blocknextai/go-packages/json"
	credentialOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/credentialoauth/domain/oauth2"
)

const (
	stateKeyPrefix    = "credentialoauth:oauth2:state:"
	stateIDRandomSize = 32
)

type StateStore struct {
	cache cache.Service
	ttl   time.Duration
}

func NewStateStore(cache cache.Service, ttl time.Duration) *StateStore {
	return &StateStore{cache: cache, ttl: ttl}
}

func (s *StateStore) Issue(ctx context.Context, payload *credentialOAuthDomainOAuth2.State) (string, error) {
	stateID, err := generateStateID()
	if err != nil {
		return "", ErrInvalidState
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return "", ErrInvalidState
	}

	if err := s.cache.Set(ctx, stateKeyPrefix+stateID, string(encoded), s.ttl); err != nil {
		return "", ErrInvalidState
	}

	return stateID, nil
}

func (s *StateStore) Consume(ctx context.Context, stateID string) (*credentialOAuthDomainOAuth2.State, error) {
	encoded, err := s.cache.GetAndDelete(ctx, stateKeyPrefix+stateID)
	if err != nil {
		return nil, ErrInvalidState
	}

	if encoded == "" {
		return nil, ErrInvalidState
	}

	var payload credentialOAuthDomainOAuth2.State
	if err := json.Unmarshal([]byte(encoded), &payload); err != nil {
		return nil, ErrInvalidState
	}

	return &payload, nil
}

func generateStateID() (string, error) {
	buf := make([]byte, stateIDRandomSize)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncode(buf), nil
}
