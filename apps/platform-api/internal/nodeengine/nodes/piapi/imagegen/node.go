package imagegen

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type PiAPIImageGenNode struct {
	nodes.Node
}

func NewPiAPIImageGenNode(nodeID string) *PiAPIImageGenNode {
	return &PiAPIImageGenNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "PiAPI Image Generation",
			Description: "Generate an image using PiAPI.",
			Icon: nodes.NodeIcon{
				Brand: "piapi",
				Glyph: "image",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"Image"},
			SubCategories: []string{"PiApi"},
			Tags: []string{
				"ai",
				"image",
				"generation",
				"media",
				"creative",
			},
			SupportedCredentials: []string{
				"piapi_api",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"model": {
						Type:        "string",
						Title:       "Model",
						Description: "PiAPI image model identifier to use.",
						Enum:        []any{"Qubico/flux1-schnell", "Qubico/flux1-dev", "Qubico/flux1-dev-advanced"},
						Default:     json.RawMessage(`"Qubico/flux1-schnell"`),
					},
					"prompt": {
						Type:        "string",
						Title:       "Prompt",
						Description: "Prompt describing the image to generate.",
					},
					"denoise": {
						Type:        "number",
						Title:       "Denoise",
						Description: "Denoising strength applied during generation.",
						Minimum:     new(0.0),
						Maximum:     new(1.0),
					},
					"negativePrompt": {
						Type:        "string",
						Title:       "Negative Prompt",
						Description: "Aspects to discourage in the generated image.",
					},
					"guidanceScale": {
						Type:        "number",
						Title:       "Guidance Scale",
						Description: "Classifier-free guidance scale for the model.",
						Minimum:     new(1.5),
						Maximum:     new(5.0),
					},
					"width": {
						Type:        "number",
						Title:       "Width",
						Description: "Width of the generated image in pixels.",
					},
					"height": {
						Type:        "number",
						Title:       "Height",
						Description: "Height of the generated image in pixels.",
					},
					"batchSize": {
						Type:        "number",
						Title:       "Batch Size",
						Description: "Number of images to generate in a single batch.",
						Minimum:     new(1.0),
						Maximum:     new(4.0),
					},
					"loraType": {
						Type:        "string",
						Title:       "Lora Type",
						Description: "Optional LoRA style to apply to the generation.",
						Enum: []any{
							"anime",
							"art",
							"disney",
							"furry",
							"mjv6",
							"realism",
							"scenery",
							"collage-artstyle",
							"creepcute",
							"cyberpunk-anime-style",
							"deco-pulse",
							"deep-sea-particle-enhencer",
							"faetastic-details",
							"fractal-geometry",
							"galactixy-illustrations-style",
							"geometric-woman",
							"graphic-portrait",
							"mat-miller-art",
							"moebius-style",
							"ob3d-isometric-3d-room",
							"paper-quilling-and-layering-style",
						},
					},
					"loraImage": {
						Type:        "string",
						Title:       "Lora Image",
						Description: "Optional reference image URL for the LoRA.",
					},
					"controlNetType": {
						Type:        "string",
						Title:       "Control Net Type",
						Description: "Type of ControlNet conditioning to apply.",
						Enum:        []any{"depth", "soft_edge", "canny"},
						Default:     json.RawMessage(`"depth"`),
					},
					"controlNetImage": {
						Type:        "string",
						Title:       "Control Net Image",
						Description: "Reference image URL used for ControlNet conditioning.",
					},
					"serviceMode": {
						Type:        "string",
						Title:       "Service Mode",
						Description: "PiAPI service mode used for the request.",
						Enum:        []any{"private", "public"},
						Default:     json.RawMessage(`"private"`),
					},
				},
				Required: []string{
					"model",
					"prompt",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"image": {
							Type:        "string",
							Description: "URL of the generated image.",
							Format:      "uri",
						},
					},
				},
			},
			HasNaturalLanguage: true,
			Annotations: nodes.NodeAnnotations{
				ReadOnly:    true,
				Destructive: new(false),
			},
		},
	}
}
