package sendgrid

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/sendgrid/sendemail"
)

func Register() {
	sendEmailNodeID := "sendgrid_send_email"
	sendEmailNode := sendemail.NewSendgridSendEmailNode(sendEmailNodeID)
	sendEmailValidator := jsonschema.New[sendemail.SendgridSendEmailExecutorInput](sendEmailNode.GetInputSchema())
	sendEmailExecutor := sendemail.NewSendgridSendEmailExecutor(sendEmailNodeID, sendEmailValidator)

	nodes.RegisterNode(sendEmailNode)
	executors.RegisterExecutor(sendEmailExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(sendEmailNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          "sendgrid",
		Name:        "SendGrid",
		Description: "Tools for sending emails via SendGrid.",
		Version:     "0.0.1",
		Tools: []nodes.NodeManager{
			sendEmailNode,
		},
	})
}
