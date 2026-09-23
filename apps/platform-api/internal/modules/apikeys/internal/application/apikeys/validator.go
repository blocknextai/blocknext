package apikeys

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/digest"
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	apiKeysDomain "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/domain/apikeys"
)

var (
	ErrAPIKeyMissing = apperror.Unauthorized("api key missing")
	ErrAPIKeyInvalid = apperror.Unauthorized("api key invalid")
	ErrAPIKeyLookup  = apperror.Internal("api key lookup failed")
)

type apiKeyValidator struct {
	repository apiKeysDomain.APIKeyRepository
}

func NewAPIKeyValidator(repository apiKeysDomain.APIKeyRepository) commonAuth.APIKeyValidator[apiKeysDomain.Scope] {
	return &apiKeyValidator{
		repository: repository,
	}
}

func (v *apiKeyValidator) Validate(ctx context.Context, rawKey string) (*commonAuth.AuthenticatedAPIKey[apiKeysDomain.Scope], error) {
	if strings.TrimSpace(rawKey) == "" {
		return nil, ErrAPIKeyMissing
	}

	apiKey, err := v.repository.GetByKeyHash(ctx, digest.SHA256Hex(rawKey))
	if err != nil {
		if errors.Is(err, apiKeysDomain.ErrAPIKeyNotFound) {
			return nil, ErrAPIKeyInvalid
		}
		return nil, ErrAPIKeyLookup.WithCause(err)
	}

	if _, updateErr := apiKey.Update(new(time.Now().UTC())); updateErr == nil {
		if err := v.repository.Update(ctx, apiKey); err != nil {
			slog.WarnContext(ctx, "failed to update api key last used at",
				"component", "apikey.validator",
				"api_key_id", apiKey.ID,
				"error", err)
		}
	}

	return &commonAuth.AuthenticatedAPIKey[apiKeysDomain.Scope]{
		ID:        apiKey.ID,
		OwnerType: apiKey.OwnerType,
		OwnerID:   apiKey.OwnerID,
		Scopes:    apiKey.Scopes,
	}, nil
}
