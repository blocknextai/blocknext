package telegram

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/adapters"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/telegram/sendmedia"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/telegram/sendmessage"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/telegram/webhook"
)

func Register() {
	nodeID := "telegram"

	sendMessageNodeID := nodeID + "_send_message"
	sendMessageNode := sendmessage.NewTelegramSendMessageNode(sendMessageNodeID)
	sendMessageValidator := jsonschema.New[sendmessage.TelegramSendMessageExecutorInput](sendMessageNode.GetInputSchema())
	sendMessageExecutor := sendmessage.NewTelegramSendMessageExecutor(sendMessageNodeID, sendMessageValidator)

	sendMediaNodeID := nodeID + "_send_media"
	sendMediaNode := sendmedia.NewTelegramSendMediaNode(sendMediaNodeID)
	sendMediaValidator := jsonschema.New[sendmedia.TelegramSendMediaExecutorInput](sendMediaNode.GetInputSchema())
	sendMediaExecutor := sendmedia.NewTelegramSendMediaExecutor(sendMediaNodeID, sendMediaValidator)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "Telegram",
			Description: "Tools for sending messages and media via Telegram.",
			Icon: mcp.ServerIcon{
				Brand: "telegram",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: sendMessageNode, Executor: sendMessageExecutor},
			{Node: sendMediaNode, Executor: sendMediaExecutor},
		},
		Adapters: []adapters.TriggerAdapter{
			webhook.NewTelegramAdapter(nodeID),
		},
	})
}
