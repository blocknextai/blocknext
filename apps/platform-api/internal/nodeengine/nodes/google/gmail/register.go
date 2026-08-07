package gmail

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/gmail/organizeemails"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/gmail/searchemails"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/gmail/sendemail"
)

func Register() {
	nodeID := "gmail"

	organizeEmailsNodeID := nodeID + "_organize_emails"
	organizeEmailsNode := organizeemails.NewGmailOrganizeEmailsNode(organizeEmailsNodeID)
	organizeEmailsValidator := jsonschema.New[organizeemails.GmailOrganizeEmailsExecutorInput](organizeEmailsNode.GetInputSchema())
	organizeEmailsExecutor := organizeemails.NewGmailOrganizeEmailsExecutor(organizeEmailsNodeID, organizeEmailsValidator)

	nodes.RegisterNode(organizeEmailsNode)
	executors.RegisterExecutor(organizeEmailsExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(organizeEmailsNode))

	searchEmailsNodeID := nodeID + "_search_emails"
	searchEmailsNode := searchemails.NewGmailSearchEmailsNode(searchEmailsNodeID)
	searchEmailsValidator := jsonschema.New[searchemails.GmailSearchEmailsExecutorInput](searchEmailsNode.GetInputSchema())
	searchEmailsExecutor := searchemails.NewGmailSearchEmailsExecutor(searchEmailsNodeID, searchEmailsValidator)

	nodes.RegisterNode(searchEmailsNode)
	executors.RegisterExecutor(searchEmailsExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(searchEmailsNode))

	sendEmailNodeID := nodeID + "_send_email"
	sendEmailNode := sendemail.NewGmailSendEmailNode(sendEmailNodeID)
	sendEmailValidator := jsonschema.New[sendemail.GmailSendEmailExecutorInput](sendEmailNode.GetInputSchema())
	sendEmailExecutor := sendemail.NewGmailSendEmailExecutor(sendEmailNodeID, sendEmailValidator)

	nodes.RegisterNode(sendEmailNode)
	executors.RegisterExecutor(sendEmailExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(sendEmailNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "Gmail",
		Description: "Tools for organizing, searching, and sending Gmail emails.",
		Icon: mcp.ServerIcon{
			Brand: "gmail",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			organizeEmailsNode,
			searchEmailsNode,
			sendEmailNode,
		},
	})
}
