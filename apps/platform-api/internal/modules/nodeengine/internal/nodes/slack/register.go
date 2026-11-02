package slack

import (
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/adapters"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/slack/sendmedia"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/slack/sendmessage"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/slack/webhook"
)

func Register(fileGateway filegateway.FileGateway) {
	nodeID := "slack"

	sendMessageNodeID := nodeID + "_send_message"
	sendMessageNode := sendmessage.NewSlackSendMessageNode(sendMessageNodeID)
	sendMessageValidator := jsonschema.New[sendmessage.SlackSendMessageExecutorInput](sendMessageNode.GetInputSchema())
	sendMessageExecutor := sendmessage.NewSlackSendMessageExecutor(sendMessageNodeID, sendMessageValidator)

	sendMediaNodeID := nodeID + "_send_media"
	sendMediaNode := sendmedia.NewSlackSendMediaNode(sendMediaNodeID)
	sendMediaValidator := jsonschema.New[sendmedia.SlackSendMediaExecutorInput](sendMediaNode.GetInputSchema())
	sendMediaExecutor := sendmedia.NewSlackSendMediaExecutor(sendMediaNodeID, sendMediaValidator, fileGateway)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "Slack",
			Description: "Tools for sending messages and media to Slack.",
			Icon: mcp.ServerIcon{
				Brand: "slack",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: sendMessageNode, Executor: sendMessageExecutor},
			{Node: sendMediaNode, Executor: sendMediaExecutor},
		},
		Adapters: []adapters.TriggerAdapter{
			webhook.NewSlackAdapter(nodeID),
		},
	})
}
