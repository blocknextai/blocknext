package searchemails

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GmailSearchEmailsNode struct {
	nodes.Node
}

func NewGmailSearchEmailsNode(nodeID string) *GmailSearchEmailsNode {
	return &GmailSearchEmailsNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Gmail Search Emails",
			Description: "Search Gmail messages using Gmail's standard search syntax and return matching email summaries.",
			Icon: nodes.NodeIcon{
				Brand: "gmail",
				Glyph: "search",
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
				"search",
				"find",
				"query",
			},
			SupportedCredentials: []string{
				"gmail_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"query": {
						Type:        "string",
						Title:       "Query",
						Description: "Gmail search query (uses Gmail's standard search syntax).",
					},
				},
				Required: []string{
					"query",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"status": {
							Type:        "boolean",
							Description: "Whether the search completed successfully.",
						},
						"result": {
							Type: "object",
							Properties: map[string]*gjs.Schema{
								"emails": {
									Type:        "array",
									Description: "List of matching email summaries.",
									Items: &gjs.Schema{
										Type: "object",
										Properties: map[string]*gjs.Schema{
											"id": {
												Type:        "string",
												Description: "Gmail message identifier.",
											},
											"threadId": {
												Type:        "string",
												Description: "Gmail thread identifier.",
											},
											"subject": {
												Type:        "string",
												Description: "Email subject.",
											},
											"sender": {
												Type:        "string",
												Description: "Email sender (From header).",
											},
											"date": {
												Type:        "string",
												Description: "Email date (Date header).",
											},
											"snippet": {
												Type:        "string",
												Description: "Short snippet of the email body.",
											},
										},
									},
								},
								"foundEmails": {
									Type:        "integer",
									Description: "Number of emails returned by the search.",
								},
								"query": {
									Type:        "string",
									Description: "The search query that was executed.",
								},
							},
						},
					},
				},
			},
			HasNaturalLanguage: true,
			Annotations: nodes.NodeAnnotations{
				ReadOnly:   true,
				Idempotent: true,
			},
		},
	}
}
