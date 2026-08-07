package discord

import (
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/discord/sendmedia"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/discord/sendmessage"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/discord/webhook"
)

func Register(fileGateway filegateway.FileGateway) {
	nodeID := "discord"

	sendMessageNodeID := nodeID + "_send_message"
	sendMessageNode := sendmessage.NewDiscordSendMessageNode(sendMessageNodeID)
	sendMessageValidator := jsonschema.New[sendmessage.DiscordSendMessageExecutorInput](sendMessageNode.GetInputSchema())
	sendMessageExecutor := sendmessage.NewDiscordSendMessageExecutor(sendMessageNodeID, sendMessageValidator)

	nodes.RegisterNode(sendMessageNode)
	executors.RegisterExecutor(sendMessageExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(sendMessageNode))

	sendMediaNodeID := nodeID + "_send_media"
	sendMediaNode := sendmedia.NewDiscordSendMediaNode(sendMediaNodeID)
	sendMediaValidator := jsonschema.New[sendmedia.DiscordSendMediaExecutorInput](sendMediaNode.GetInputSchema())
	sendMediaExecutor := sendmedia.NewDiscordSendMediaExecutor(sendMediaNodeID, sendMediaValidator, fileGateway)

	nodes.RegisterNode(sendMediaNode)
	executors.RegisterExecutor(sendMediaExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(sendMediaNode))

	adapters.RegisterAdapter(webhook.NewDiscordAdapter(nodeID))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "Discord",
		Description: "Tools for sending messages and media to Discord.",
		Icon: mcp.ServerIcon{
			Brand: "discord",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			sendMessageNode,
			sendMediaNode,
		},
	})
}
