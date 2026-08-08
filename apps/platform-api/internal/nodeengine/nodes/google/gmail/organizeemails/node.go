package organizeemails

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GmailOrganizeEmailsNode struct {
	nodes.Node
}

func NewGmailOrganizeEmailsNode(nodeID string) *GmailOrganizeEmailsNode {
	return &GmailOrganizeEmailsNode{
		Node: nodes.Node{
			ID:          nodeID,
			Kind:        nodes.NodeKindAction,
			Version:     "0.0.1",
			Name:        "Gmail Organize Emails",
			Description: "Organize Gmail messages by applying labels, stars, categories, read state, or trash actions to messages matching a search query.",
			Icon: nodes.NodeIcon{
				Brand: "gmail",
				Glyph: "organize",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"Mailing"},
			SubCategories: []string{"Gmail"},
			Tags: []string{
				"google",
				"email",
				"mail",
				"organize",
				"label",
				"automation",
			},
			SupportedCredentials: []string{
				"gmail_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"keywords": {
						Type:        "string",
						Title:       "Keywords",
						Description: "Gmail search query used to select emails to organize.",
					},
					"actions": {
						Type:        "array",
						Title:       "Actions",
						Description: "Actions to apply to the matching emails.",
						Items: &gjs.Schema{
							Type: "object",
							Properties: map[string]*gjs.Schema{
								"type": {
									Type:        "string",
									Description: "Action type: label, star, read, category, or trash.",
									Enum:        []any{"label", "star", "read", "category", "trash"},
								},
								"labelName": {
									Type:        "string",
									Description: "Label name to apply when type is \"label\".",
								},
								"category": {
									Type:        "string",
									Description: "Gmail category when type is \"category\".",
									Enum:        []any{"primary", "social", "promotions", "updates", "forums"},
								},
								"star": {
									Type:        "boolean",
									Description: "Whether to star (true) or unstar (false) when type is \"star\".",
								},
								"read": {
									Type:        "boolean",
									Description: "Whether to mark as read (true) or unread (false) when type is \"read\".",
								},
								"trash": {
									Type:        "boolean",
									Description: "Whether to move emails to trash when type is \"trash\".",
								},
							},
						},
					},
				},
				Required: []string{
					"keywords",
					"actions",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"status": {
							Type:        "boolean",
							Description: "Whether the organize operation completed successfully.",
						},
						"result": {
							Type: "object",
							Properties: map[string]*gjs.Schema{
								"processedEmails": {
									Type:        "integer",
									Description: "Number of emails successfully processed.",
								},
								"appliedActions": {
									Type:        "array",
									Description: "Human-readable list of actions that were applied.",
									Items:       &gjs.Schema{Type: "string"},
								},
								"scope": {
									Type:        "string",
									Description: "Scope of the operation (search query or \"all_emails\").",
								},
							},
						},
					},
				},
			},
			HasNaturalLanguage: true,
			Annotations: nodes.NodeAnnotations{
				Destructive: new(true),
			},
		},
	}
}
