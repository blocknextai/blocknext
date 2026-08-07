package slack

import (
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/slack/sendmedia"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/slack/sendmessage"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/slack/webhook"
)

func Register(fileGateway filegateway.FileGateway) {
	nodeID := "slack"

	sendMessageNodeID := nodeID + "_send_message"
	sendMessageNode := sendmessage.NewSlackSendMessageNode(sendMessageNodeID)
	sendMessageValidator := jsonschema.New[sendmessage.SlackSendMessageExecutorInput](sendMessageNode.GetInputSchema())
	sendMessageExecutor := sendmessage.NewSlackSendMessageExecutor(sendMessageNodeID, sendMessageValidator)

	nodes.RegisterNode(sendMessageNode)
	executors.RegisterExecutor(sendMessageExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(sendMessageNode))

	sendMediaNodeID := nodeID + "_send_media"
	sendMediaNode := sendmedia.NewSlackSendMediaNode(sendMediaNodeID)
	sendMediaValidator := jsonschema.New[sendmedia.SlackSendMediaExecutorInput](sendMediaNode.GetInputSchema())
	sendMediaExecutor := sendmedia.NewSlackSendMediaExecutor(sendMediaNodeID, sendMediaValidator, fileGateway)

	nodes.RegisterNode(sendMediaNode)
	executors.RegisterExecutor(sendMediaExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(sendMediaNode))

	adapters.RegisterAdapter(webhook.NewSlackAdapter(nodeID))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "Slack",
		Description: "Tools for sending messages and media to Slack.",
		Icon: mcp.ServerIcon{
			Brand: "slack",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			sendMessageNode,
			sendMediaNode,
		},
	})
}
