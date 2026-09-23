package auth

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/pkce"
	"github.com/blocknextai/go-packages/uuid"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	accountDomainUserNonces "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/usernonces"
)

type CreateUserNonceCommand struct {
	AuthProvider domain.AuthProvider
	ProviderID   *string
}

var (
	errFailedToGenerateOAuthURL = apperror.Internal("failed to generate oauth url")
)

type CreateUserNonceResponse struct {
	Nonce        string `json:"nonce"`
	URL          string `json:"url"`
	LoginMessage string `json:"loginMessage,omitempty"`
}

func (s *Service) CreateUserNonce(ctx context.Context, command *CreateUserNonceCommand) (*CreateUserNonceResponse, error) {
	var response *CreateUserNonceResponse

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		challenge, err := pkce.GenerateChallenge()
		if err != nil {
			return err
		}

		nonce := uuid.NewV7().String()

		userNonce, err := accountDomainUserNonces.NewUserNonce(
			command.AuthProvider,
			command.ProviderID,
			nonce,
			challenge.CodeVerifier,
			challenge.CodeChallenge,
			challenge.CodeChallengeMethod,
		)
		if err != nil {
			return err
		}

		err = s.userNonceRepository.Create(txCtx, userNonce)
		if err != nil {
			return err
		}

		authProvider, err := s.authProviderRegistry.GetProvider(command.AuthProvider)
		if err != nil {
			return err
		}

		oauthURL, err := authProvider.GenerateOAuthURL(userNonce)
		if err != nil {
			return errFailedToGenerateOAuthURL.WithCause(err)
		}

		loginMessage := authProvider.BuildLoginMessage(nonce)

		response = &CreateUserNonceResponse{
			Nonce:        nonce,
			URL:          oauthURL,
			LoginMessage: loginMessage,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
