package createusernonce

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/pkce"
	"github.com/blocknextai/go-packages/uuid"
	"github.com/blocknextai/platform-api/internal/account/application/auth/createusertoken"
	accountDomainUserNonces "github.com/blocknextai/platform-api/internal/account/domain/usernonces"
)

type Handler struct {
	userNonceRepository  accountDomainUserNonces.UserNonceRepository
	authProviderRegistry createusertoken.AuthProviderRegistry
	transactionManager   database.TransactionManager
}

func New(
	userNonceRepository accountDomainUserNonces.UserNonceRepository,
	authProviderRegistry createusertoken.AuthProviderRegistry,
	transactionManager database.TransactionManager,
) *Handler {
	return &Handler{
		userNonceRepository:  userNonceRepository,
		authProviderRegistry: authProviderRegistry,
		transactionManager:   transactionManager,
	}
}

func (h *Handler) Handle(ctx context.Context, command *CreateUserNonceCommand) (*CreateUserNonceResponse, error) {
	var response *CreateUserNonceResponse

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
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

		err = h.userNonceRepository.Create(txCtx, userNonce)
		if err != nil {
			return err
		}

		authProvider, err := h.authProviderRegistry.GetProvider(command.AuthProvider)
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
