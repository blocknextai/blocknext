package servers

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrServerNotFound    = apperror.NotFound("mcp server not found")
	ErrDuplicateServerID = apperror.Internal("duplicate mcp server id")
)
