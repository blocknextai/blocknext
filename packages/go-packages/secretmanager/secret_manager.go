package secretmanager

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/blocknextai/go-packages/base64"
	"github.com/blocknextai/go-packages/hex"
	"github.com/blocknextai/go-packages/json"
)

type secretManager struct {
	secretKey []byte
}

// New returns a SecretManager that uses the given hex-encoded 32-byte secret
// key for AES-256-GCM encryption. It returns ErrInvalidKeyFormat if the key
// is not valid hexadecimal or ErrInvalidKeyLength if it does not decode to 32
// bytes.
func New(
	secretKey string,
) (SecretManager, error) {
	key, err := hex.Decode(secretKey)
	if err != nil {
		return nil, ErrInvalidKeyFormat
	}

	if len(key) != 32 {
		return nil, ErrInvalidKeyLength
	}

	return &secretManager{
		secretKey: key,
	}, nil
}

func (s *secretManager) Encrypt(value any) (string, error) {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.secretKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, jsonData, nil)

	return base64.Encode(ciphertext), nil
}

func (s *secretManager) Decrypt(encryptedText string, output any) error {
	ciphertext, err := base64.Decode(encryptedText)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(s.secretKey)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return ErrCiphertextTooShort
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(plaintext, output); err != nil {
		return err
	}

	return nil
}
