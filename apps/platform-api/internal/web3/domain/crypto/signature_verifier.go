package crypto

type SignatureVerifier interface {
	VerifySignature(address string, nonce string, signature string) error
	PersonalEcRecover(message, sigHex string) (string, error)
	BuildLoginMessage(nonce string) string
}
