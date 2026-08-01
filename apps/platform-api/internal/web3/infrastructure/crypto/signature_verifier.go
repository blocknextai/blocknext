package crypto

import (
	"strings"

	"github.com/blocknextai/go-packages/cast"
	web3DomainCrypto "github.com/blocknextai/platform-api/internal/web3/domain/crypto"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type ethereumSignatureVerifier struct {
	loginMessage string
}

func NewEthereumSignatureVerifier(loginMessage string) web3DomainCrypto.SignatureVerifier {
	return &ethereumSignatureVerifier{
		loginMessage: loginMessage,
	}
}

func (v *ethereumSignatureVerifier) BuildLoginMessage(nonce string) string {
	return strings.ReplaceAll(v.loginMessage, "{nonce}", nonce)
}

func (v *ethereumSignatureVerifier) VerifySignature(address string, nonce string, signature string) error {
	message := v.BuildLoginMessage(nonce)

	recoveredAddr, err := v.PersonalEcRecover(message, signature)
	if err != nil || !strings.EqualFold(recoveredAddr, address) {
		return web3DomainCrypto.ErrInvalidSignature
	}

	return nil
}

func (v *ethereumSignatureVerifier) PersonalEcRecover(message, sigHex string) (string, error) {
	var builder strings.Builder
	builder.WriteString("\x19Ethereum Signed Message:\n")
	builder.WriteString(cast.ToString(len(message)))
	builder.WriteString(message)
	msg := []byte(builder.String())
	msgHash := crypto.Keccak256Hash(msg)

	sigBytes, err := hexutil.Decode(sigHex)
	if err != nil {
		return "", err
	}
	if sigBytes[crypto.RecoveryIDOffset] == 27 || sigBytes[crypto.RecoveryIDOffset] == 28 {
		sigBytes[crypto.RecoveryIDOffset] -= 27
	}

	pubKey, err := crypto.SigToPub(msgHash.Bytes(), sigBytes)
	if err != nil {
		return "", err
	}

	addr := crypto.PubkeyToAddress(*pubKey).Hex()
	return addr, nil
}
