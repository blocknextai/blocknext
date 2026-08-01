// Package bcrypt provides a bcrypt-based implementation of the hashing.Hasher
// interface.
package bcrypt

import (
	"github.com/blocknextai/go-packages/hashing"
	"golang.org/x/crypto/bcrypt"
)

type bcryptHasher struct {
	cost int
}

// New returns a hashing.Hasher that uses bcrypt with the given cost.
func New(cost int) hashing.Hasher {
	return &bcryptHasher{cost: cost}
}

func (b *bcryptHasher) Generate(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (b *bcryptHasher) Compare(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
