package adapters

import (
	"strings"
)

type VerificationRequest struct {
	Method      string
	Headers     map[string][]string
	Body        []byte
	QueryParams map[string]string
}

func (r *VerificationRequest) Header(name string) string {
	for key, values := range r.Headers {
		if strings.EqualFold(key, name) && len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

type VerificationResponse struct {
	Body       []byte
	StatusCode int
}

type WebhookVerifier interface {
	Verify(request *VerificationRequest, secret *string) (*VerificationResponse, error)
	ValidateSignature(request *VerificationRequest, secret *string) error
}

func AsWebhookVerifier(adapter TriggerAdapter) (WebhookVerifier, bool) {
	if adapter == nil {
		return nil, false
	}
	v, ok := adapter.(WebhookVerifier)
	return v, ok
}
