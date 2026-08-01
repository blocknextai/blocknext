package publishpost

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type LinkedinPublishPostNode struct {
	nodes.Node
}

func NewLinkedinPublishPostNode(nodeID string) *LinkedinPublishPostNode {
	return &LinkedinPublishPostNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "LinkedIn Publish Post",
			Description: "Publish a post to a LinkedIn account.",
			Icon: nodes.NodeIcon{
				Light: nodeID,
				Dark:  nodeID,
			},
			Categories:    []string{"Publishing"},
			SubCategories: []string{"LinkedIn"},
			Tags: []string{
				"social",
				"social-media",
				"professional",
				"b2b",
				"publish",
				"post",
				"share",
			},
			SupportedCredentials: []string{
				"linkedin_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"text": {
						Type:        "string",
						Title:       "Text",
						Description: "Text content of the LinkedIn post.",
					},
					"visibility": {
						Type:        "string",
						Title:       "Visibility",
						Description: "Visibility of the post (PUBLIC or CONNECTIONS).",
						Enum:        []any{"PUBLIC", "CONNECTIONS"},
						Default:     json.RawMessage(`"PUBLIC"`),
					},
				},
				Required: []string{
					"text",
					"visibility",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"status": {
							Type:        "boolean",
							Description: "Whether the post was published successfully.",
						},
						"postId": {
							Type:        "string",
							Description: "Identifier of the published LinkedIn post.",
						},
						"postUrl": {
							Type:        "string",
							Description: "URL of the published LinkedIn post.",
							Format:      "uri",
						},
					},
				},
			},
			HasNaturalLanguage: true,
			Annotations: nodes.NodeAnnotations{
				Destructive: new(false),
			},
		},
	}
}
