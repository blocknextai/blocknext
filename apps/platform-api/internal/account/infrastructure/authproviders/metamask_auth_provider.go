package authproviders

import (
	"context"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/platform-api/internal/account/application/auth/createusertoken"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	accountDomainUserNonces "github.com/blocknextai/platform-api/internal/account/domain/usernonces"
	web3DomainCrypto "github.com/blocknextai/platform-api/internal/web3/domain/crypto"
)

var (
	ErrMetamaskSignatureVerificationFailed = apperror.Unauthorized("metamask signature verification failed")
	ErrMetamaskInvalidNonce                = apperror.Validation("metamask invalid nonce")
	ErrMetamaskInvalidWalletAddress        = apperror.Validation("metamask invalid wallet address")
	ErrMetamaskInvalidSignature            = apperror.Validation("metamask invalid signature")
)

type MetamaskAuthProvider struct {
	signatureVerifier   web3DomainCrypto.SignatureVerifier
	userNonceRepository accountDomainUserNonces.UserNonceRepository
}

func NewMetamaskAuthProvider(
	signatureVerifier web3DomainCrypto.SignatureVerifier,
	userNonceRepository accountDomainUserNonces.UserNonceRepository,
) createusertoken.AuthProvider {
	return &MetamaskAuthProvider{
		signatureVerifier:   signatureVerifier,
		userNonceRepository: userNonceRepository,
	}
}

func (s *MetamaskAuthProvider) Validate(ctx context.Context, request createusertoken.Request) (*createusertoken.Response, error) {
	nonceStr, ok := request.Payload["nonce"].(string)
	if !ok {
		return nil, ErrMetamaskInvalidNonce
	}

	nonce, err := s.userNonceRepository.GetByNonceAndAuthProvider(ctx, nonceStr, accountDomain.AuthProviderMetamask)
	if err != nil {
		return nil, err
	}

	walletAddress, ok := request.Payload["walletAddress"].(string)
	if !ok {
		return nil, ErrMetamaskInvalidWalletAddress
	}
	signature, ok := request.Payload["signature"].(string)
	if !ok {
		return nil, ErrMetamaskInvalidSignature
	}

	err = s.signatureVerifier.VerifySignature(walletAddress, nonce.Nonce, signature)
	if err != nil {
		return nil, ErrMetamaskSignatureVerificationFailed
	}

	return &createusertoken.Response{
		ProviderID:  walletAddress,
		Identifier:  walletAddress,
		DisplayName: walletAddress,
		Nonce:       nonce,
	}, nil
}

func (s *MetamaskAuthProvider) GenerateOAuthURL(userNonce *accountDomainUserNonces.UserNonce) (string, error) {
	return "", nil
}

func (s *MetamaskAuthProvider) BuildLoginMessage(nonce string) string {
	return s.signatureVerifier.BuildLoginMessage(nonce)
}
