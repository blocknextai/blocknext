package telegram

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/telegram/sendmedia"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/telegram/sendmessage"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/telegram/webhook"
)

func Register() {
	nodeID := "telegram"

	sendMessageNodeID := nodeID + "_send_message"
	sendMessageNode := sendmessage.NewTelegramSendMessageNode(sendMessageNodeID)
	sendMessageValidator := jsonschema.New[sendmessage.TelegramSendMessageExecutorInput](sendMessageNode.GetInputSchema())
	sendMessageExecutor := sendmessage.NewTelegramSendMessageExecutor(sendMessageNodeID, sendMessageValidator)

	nodes.RegisterNode(sendMessageNode)
	executors.RegisterExecutor(sendMessageExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(sendMessageNode))

	sendMediaNodeID := nodeID + "_send_media"
	sendMediaNode := sendmedia.NewTelegramSendMediaNode(sendMediaNodeID)
	sendMediaValidator := jsonschema.New[sendmedia.TelegramSendMediaExecutorInput](sendMediaNode.GetInputSchema())
	sendMediaExecutor := sendmedia.NewTelegramSendMediaExecutor(sendMediaNodeID, sendMediaValidator)

	nodes.RegisterNode(sendMediaNode)
	executors.RegisterExecutor(sendMediaExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(sendMediaNode))

	adapters.RegisterAdapter(webhook.NewTelegramAdapter(nodeID))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "Telegram",
		Description: "Tools for sending messages and media via Telegram.",
		Icon: mcp.ServerIcon{
			Brand: "telegram",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			sendMessageNode,
			sendMediaNode,
		},
	})
}
