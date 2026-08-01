package triggers

import (
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrTriggerNotFound    = apperror.NotFound("trigger not found")
	ErrNotWebhookTrigger  = apperror.Validation("not a webhook trigger")
	ErrNotScheduleTrigger = apperror.Validation("not a schedule trigger")
	ErrUnsupportedType    = apperror.Validation("unsupported trigger type")
	ErrWebhookTokenTaken  = apperror.Conflict("webhook token already taken")
)
