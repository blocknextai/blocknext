package gmail

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/gmail/organizeemails"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/gmail/searchemails"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/gmail/sendemail"
)

func Register() {
	nodeID := "gmail"

	organizeEmailsNodeID := nodeID + "_organize_emails"
	organizeEmailsNode := organizeemails.NewGmailOrganizeEmailsNode(organizeEmailsNodeID)
	organizeEmailsValidator := jsonschema.New[organizeemails.GmailOrganizeEmailsExecutorInput](organizeEmailsNode.GetInputSchema())
	organizeEmailsExecutor := organizeemails.NewGmailOrganizeEmailsExecutor(organizeEmailsNodeID, organizeEmailsValidator)

	searchEmailsNodeID := nodeID + "_search_emails"
	searchEmailsNode := searchemails.NewGmailSearchEmailsNode(searchEmailsNodeID)
	searchEmailsValidator := jsonschema.New[searchemails.GmailSearchEmailsExecutorInput](searchEmailsNode.GetInputSchema())
	searchEmailsExecutor := searchemails.NewGmailSearchEmailsExecutor(searchEmailsNodeID, searchEmailsValidator)

	sendEmailNodeID := nodeID + "_send_email"
	sendEmailNode := sendemail.NewGmailSendEmailNode(sendEmailNodeID)
	sendEmailValidator := jsonschema.New[sendemail.GmailSendEmailExecutorInput](sendEmailNode.GetInputSchema())
	sendEmailExecutor := sendemail.NewGmailSendEmailExecutor(sendEmailNodeID, sendEmailValidator)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "Gmail",
			Description: "Tools for organizing, searching, and sending Gmail emails.",
			Icon: mcp.ServerIcon{
				Brand: "gmail",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: organizeEmailsNode, Executor: organizeEmailsExecutor},
			{Node: searchEmailsNode, Executor: searchEmailsExecutor},
			{Node: sendEmailNode, Executor: sendEmailExecutor},
		},
	})
}
