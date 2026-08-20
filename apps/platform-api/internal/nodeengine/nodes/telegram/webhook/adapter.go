package webhook

import (
	"github.com/blocknextai/go-packages/json"
	nodeEngineDomainAdapters "github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
)

type TelegramAdapter struct {
	nodeEngineDomainAdapters.Adapter
}

func NewTelegramAdapter(nodeID string) *TelegramAdapter {
	return &TelegramAdapter{
		ID:   nodeID,
		Name: "Telegram",
	}
}

type telegramPayload struct {
	Message struct {
		Text string `json:"text"`
		From struct {
			Username  string `json:"username"`
			FirstName string `json:"first_name"`
		} `json:"from"`
	} `json:"message"`
}

func (a *TelegramAdapter) Adapt(raw map[string]any) (*nodeEngineDomainAdapters.TriggerContext, error) {
	var payload telegramPayload
	if err := json.ArgsToStruct(raw, &payload); err != nil {
		return nil, err
	}

	sender := payload.Message.From.Username
	if sender == "" {
		sender = payload.Message.From.FirstName
	}

	return &nodeEngineDomainAdapters.TriggerContext{
		Source:  a.ID,
		Sender:  sender,
		Prompt:  payload.Message.Text,
		Payload: raw,
	}, nil
}
