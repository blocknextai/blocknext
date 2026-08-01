package passwordpolicy

import (
	"context"
)

type Policy interface {
	Check(ctx context.Context, password string, userInputs []string) error
}
