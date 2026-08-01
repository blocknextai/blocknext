package common

import (
	pkgEmail "github.com/blocknextai/go-packages/email"
	"github.com/blocknextai/go-packages/hashing"
	commonInfrastructure "github.com/blocknextai/platform-api/internal/common/infrastructure"
	"github.com/blocknextai/platform-api/internal/config"
)

type Dependencies struct {
	EmailSenderOptions config.EmailSenderOptions
	BcryptCost         int
}

type Module struct {
	EmailSender    pkgEmail.EmailSender
	PasswordHasher hashing.Hasher
}

func NewModule(deps Dependencies) *Module {
	return &Module{
		EmailSender:    commonInfrastructure.NewEmailSender(deps.EmailSenderOptions),
		PasswordHasher: commonInfrastructure.NewPasswordHasher(deps.BcryptCost),
	}
}
