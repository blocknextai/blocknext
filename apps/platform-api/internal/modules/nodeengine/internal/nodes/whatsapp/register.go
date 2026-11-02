package whatsapp

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/adapters"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/whatsapp/sendmedia"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/whatsapp/sendtemplate"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/whatsapp/sendtextmessage"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/whatsapp/webhook"
)

func Register() {
	nodeID := "whatsapp"

	sendMessageNodeID := nodeID + "_send_text_message"
	sendMessageNode := sendtextmessage.NewWhatsAppSendTextMessageNode(sendMessageNodeID)
	sendMessageValidator := jsonschema.New[sendtextmessage.WhatsAppSendTextMessageExecutorInput](sendMessageNode.GetInputSchema())
	sendMessageExecutor := sendtextmessage.NewWhatsAppSendTextMessageExecutor(sendMessageNodeID, sendMessageValidator)

	sendTemplateNodeID := nodeID + "_send_template"
	sendTemplateNode := sendtemplate.NewWhatsAppSendTemplateNode(sendTemplateNodeID)
	sendTemplateValidator := jsonschema.New[sendtemplate.WhatsAppSendTemplateExecutorInput](sendTemplateNode.GetInputSchema())
	sendTemplateExecutor := sendtemplate.NewWhatsAppSendTemplateExecutor(sendTemplateNodeID, sendTemplateValidator)

	sendMediaNodeID := nodeID + "_send_media"
	sendMediaNode := sendmedia.NewWhatsAppSendMediaNode(sendMediaNodeID)
	sendMediaValidator := jsonschema.New[sendmedia.WhatsAppSendMediaExecutorInput](sendMediaNode.GetInputSchema())
	sendMediaExecutor := sendmedia.NewWhatsAppSendMediaExecutor(sendMediaNodeID, sendMediaValidator)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "WhatsApp",
			Description: "Tools for sending text messages, templates, and media via WhatsApp.",
			Icon: mcp.ServerIcon{
				Brand: "whatsapp",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: sendMessageNode, Executor: sendMessageExecutor},
			{Node: sendTemplateNode, Executor: sendTemplateExecutor},
			{Node: sendMediaNode, Executor: sendMediaExecutor},
		},
		Adapters: []adapters.TriggerAdapter{
			webhook.NewWhatsAppAdapter(nodeID),
		},
	})
}
