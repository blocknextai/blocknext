package whatsapp

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/adapters"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/whatsapp/sendmedia"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/whatsapp/sendtemplate"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/whatsapp/sendtextmessage"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/whatsapp/webhook"
)

func Register() {
	nodeID := "whatsapp"

	sendMessageNodeID := nodeID + "_send_text_message"
	sendMessageNode := sendtextmessage.NewWhatsAppSendTextMessageNode(sendMessageNodeID)
	sendMessageValidator := jsonschema.New[sendtextmessage.WhatsAppSendTextMessageExecutorInput](sendMessageNode.GetInputSchema())
	sendMessageExecutor := sendtextmessage.NewWhatsAppSendTextMessageExecutor(sendMessageNodeID, sendMessageValidator)

	nodes.RegisterNode(sendMessageNode)
	executors.RegisterExecutor(sendMessageExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(sendMessageNode))

	sendTemplateNodeID := nodeID + "_send_template"
	sendTemplateNode := sendtemplate.NewWhatsAppSendTemplateNode(sendTemplateNodeID)
	sendTemplateValidator := jsonschema.New[sendtemplate.WhatsAppSendTemplateExecutorInput](sendTemplateNode.GetInputSchema())
	sendTemplateExecutor := sendtemplate.NewWhatsAppSendTemplateExecutor(sendTemplateNodeID, sendTemplateValidator)

	nodes.RegisterNode(sendTemplateNode)
	executors.RegisterExecutor(sendTemplateExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(sendTemplateNode))

	sendMediaNodeID := nodeID + "_send_media"
	sendMediaNode := sendmedia.NewWhatsAppSendMediaNode(sendMediaNodeID)
	sendMediaValidator := jsonschema.New[sendmedia.WhatsAppSendMediaExecutorInput](sendMediaNode.GetInputSchema())
	sendMediaExecutor := sendmedia.NewWhatsAppSendMediaExecutor(sendMediaNodeID, sendMediaValidator)

	nodes.RegisterNode(sendMediaNode)
	executors.RegisterExecutor(sendMediaExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(sendMediaNode))

	adapters.RegisterAdapter(webhook.NewWhatsAppAdapter(nodeID))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "WhatsApp",
		Description: "Tools for sending text messages, templates, and media via WhatsApp.",
		Version:     "0.0.1",
		Tools: []nodes.NodeManager{
			sendMessageNode,
			sendTemplateNode,
			sendMediaNode,
		},
	})
}
