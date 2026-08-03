package webhook

import (
	"github.com/blocknextai/go-packages/json"
	nodeEngineDomainAdapters "github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
)

type DiscordAdapter struct {
	nodeEngineDomainAdapters.Adapter
}

func NewDiscordAdapter(nodeID string) *DiscordAdapter {
	return &DiscordAdapter{
		Adapter: nodeEngineDomainAdapters.Adapter{
			ID:   nodeID,
			Name: "Discord",
		},
	}
}

type discordPayload struct {
	Content string `json:"content"`
	Author  struct {
		Username string `json:"username"`
	} `json:"author"`
}

func (a *DiscordAdapter) Adapt(raw map[string]any) (*nodeEngineDomainAdapters.TriggerContext, error) {
	var payload discordPayload
	if err := json.ArgsToStruct(raw, &payload); err != nil {
		return nil, err
	}

	return &nodeEngineDomainAdapters.TriggerContext{
		Source:  a.ID,
		Sender:  payload.Author.Username,
		Prompt:  payload.Content,
		Payload: raw,
	}, nil
}
