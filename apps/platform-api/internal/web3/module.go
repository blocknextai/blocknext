package web3

import (
	web3DomainCrypto "github.com/blocknextai/platform-api/internal/web3/domain/crypto"
	web3InfrastructureCrypto "github.com/blocknextai/platform-api/internal/web3/infrastructure/crypto"
)

type Dependencies struct {
	LoginMessage string
}

type Module struct {
	SignatureVerifier web3DomainCrypto.SignatureVerifier
}

func NewModule(deps Dependencies) (*Module, error) {
	signatureVerifier := web3InfrastructureCrypto.NewEthereumSignatureVerifier(deps.LoginMessage)

	return &Module{
		SignatureVerifier: signatureVerifier,
	}, nil
}
