package sendgrid

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/sendgrid/sendemail"
)

func Register() {
	sendEmailNodeID := "sendgrid_send_email"
	sendEmailNode := sendemail.NewSendgridSendEmailNode(sendEmailNodeID)
	sendEmailValidator := jsonschema.New[sendemail.SendgridSendEmailExecutorInput](sendEmailNode.GetInputSchema())
	sendEmailExecutor := sendemail.NewSendgridSendEmailExecutor(sendEmailNodeID, sendEmailValidator)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          "sendgrid",
			Name:        "SendGrid",
			Description: "Tools for sending emails via SendGrid.",
			Icon: mcp.ServerIcon{
				Brand: "sendgrid",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: sendEmailNode, Executor: sendEmailExecutor},
		},
	})
}
