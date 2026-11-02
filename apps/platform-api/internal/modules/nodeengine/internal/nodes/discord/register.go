package discord

import (
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/adapters"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/discord/sendmedia"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/discord/sendmessage"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/discord/webhook"
)

func Register(fileGateway filegateway.FileGateway) {
	nodeID := "discord"

	sendMessageNodeID := nodeID + "_send_message"
	sendMessageNode := sendmessage.NewDiscordSendMessageNode(sendMessageNodeID)
	sendMessageValidator := jsonschema.New[sendmessage.DiscordSendMessageExecutorInput](sendMessageNode.GetInputSchema())
	sendMessageExecutor := sendmessage.NewDiscordSendMessageExecutor(sendMessageNodeID, sendMessageValidator)

	sendMediaNodeID := nodeID + "_send_media"
	sendMediaNode := sendmedia.NewDiscordSendMediaNode(sendMediaNodeID)
	sendMediaValidator := jsonschema.New[sendmedia.DiscordSendMediaExecutorInput](sendMediaNode.GetInputSchema())
	sendMediaExecutor := sendmedia.NewDiscordSendMediaExecutor(sendMediaNodeID, sendMediaValidator, fileGateway)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "Discord",
			Description: "Tools for sending messages and media to Discord.",
			Icon: mcp.ServerIcon{
				Brand: "discord",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: sendMessageNode, Executor: sendMessageExecutor},
			{Node: sendMediaNode, Executor: sendMediaExecutor},
		},
		Adapters: []adapters.TriggerAdapter{
			webhook.NewDiscordAdapter(nodeID),
		},
	})
}
