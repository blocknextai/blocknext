package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"strings"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/hex"
	"github.com/blocknextai/go-packages/json"
	nodeEngineDomainAdapters "github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
)

type WhatsAppAdapter struct {
	nodeEngineDomainAdapters.Adapter
}

func NewWhatsAppAdapter(nodeID string) *WhatsAppAdapter {
	return &WhatsAppAdapter{
		ID:   nodeID,
		Name: "WhatsApp",
	}
}

type whatsappPayload struct {
	Entry []struct {
		Changes []struct {
			Value struct {
				Messages []struct {
					From string `json:"from"`
					Text struct {
						Body string `json:"body"`
					} `json:"text"`
				} `json:"messages"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

var (
	ErrInvalidWhatsAppPayload    = apperror.Internal("invalid payload")
	ErrWhatsAppVerificationToken = apperror.Unauthorized("invalid whatsapp verify_token")
	ErrWhatsAppMissingSignature  = apperror.Unauthorized("missing whatsapp signature header")
	ErrWhatsAppInvalidSignature  = apperror.Unauthorized("invalid whatsapp signature")
)

func (a *WhatsAppAdapter) Adapt(raw map[string]any) (*nodeEngineDomainAdapters.TriggerContext, error) {
	var payload whatsappPayload
	if err := json.ArgsToStruct(raw, &payload); err != nil {
		return nil, err
	}

	if len(payload.Entry) == 0 || len(payload.Entry[0].Changes) == 0 || len(payload.Entry[0].Changes[0].Value.Messages) == 0 {
		return nil, ErrInvalidWhatsAppPayload
	}

	message := payload.Entry[0].Changes[0].Value.Messages[0]

	return &nodeEngineDomainAdapters.TriggerContext{
		Source:  a.ID,
		Sender:  message.From,
		Prompt:  message.Text.Body,
		Payload: raw,
	}, nil
}

func (a *WhatsAppAdapter) Verify(request *nodeEngineDomainAdapters.VerificationRequest, secret *string) (*nodeEngineDomainAdapters.VerificationResponse, error) {
	if request.Method != "GET" {
		return nil, nil
	}

	mode := request.QueryParams["hub.mode"]
	challenge := request.QueryParams["hub.challenge"]
	verifyToken := request.QueryParams["hub.verify_token"]

	if mode != "subscribe" {
		return nil, nil
	}

	if secret == nil || strings.TrimSpace(*secret) == "" || verifyToken != *secret {
		return nil, ErrWhatsAppVerificationToken
	}

	return &nodeEngineDomainAdapters.VerificationResponse{
		Body:       []byte(challenge),
		StatusCode: 200,
	}, nil
}

func (a *WhatsAppAdapter) ValidateSignature(request *nodeEngineDomainAdapters.VerificationRequest, secret *string) error {
	if secret == nil || strings.TrimSpace(*secret) == "" {
		return nil
	}

	signature := request.Header("X-Hub-Signature-256")
	if strings.TrimSpace(signature) == "" {
		return ErrWhatsAppMissingSignature
	}

	if !strings.HasPrefix(signature, "sha256=") {
		return ErrWhatsAppInvalidSignature
	}

	mac := hmac.New(sha256.New, []byte(*secret))
	mac.Write(request.Body)
	expected := "sha256=" + hex.Encode(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return ErrWhatsAppInvalidSignature
	}

	return nil
}
