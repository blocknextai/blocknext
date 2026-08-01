package translate

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type DeeplTranslateNode struct {
	nodes.Node
}

func NewDeeplTranslateNode(nodeID string) *DeeplTranslateNode {
	return &DeeplTranslateNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "DeepL Translate",
			Description: "Translate text using DeepL.",
			Icon: nodes.NodeIcon{
				Light: nodeID,
				Dark:  nodeID,
			},
			Categories:    []string{"AI"},
			SubCategories: []string{"DeepL"},
			Tags: []string{
				"translation",
				"translate",
				"language",
				"localization",
				"text",
				"i18n",
			},
			SupportedCredentials: []string{
				"deepl_api",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"text": {
						Type:        "string",
						Title:       "Text",
						Description: "Text to translate.",
					},
					"sourceLanguage": {
						Type:        "string",
						Title:       "Source Language",
						Description: "Source language code, or 'auto' to auto-detect.",
						Default:     json.RawMessage(`"auto"`),
					},
					"targetLanguage": {
						Type:        "string",
						Title:       "Target Language",
						Description: "Target language code (e.g. en-US, de, fr).",
						Default:     json.RawMessage(`"en-US"`),
					},
				},
				Required: []string{
					"text",
					"sourceLanguage",
					"targetLanguage",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"text": {
							Type:        "string",
							Description: "Translated text.",
						},
					},
				},
			},
			HasNaturalLanguage: true,
			Annotations: nodes.NodeAnnotations{
				ReadOnly:    true,
				Idempotent:  true,
				Destructive: new(false),
			},
		},
	}
}
