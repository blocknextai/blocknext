package videogen

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type PiAPIVideoGenNode struct {
	nodes.Node
}

func NewPiAPIVideoGenNode(nodeID string) *PiAPIVideoGenNode {
	return &PiAPIVideoGenNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "PiAPI Video Generation",
			Description: "Generate a video using PiAPI.",
			Icon: nodes.NodeIcon{
				Light: nodeID,
				Dark:  nodeID,
			},
			Categories:    []string{"Video"},
			SubCategories: []string{"PiApi"},
			Tags: []string{
				"ai",
				"video",
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
					"prompt": {
						Type:        "string",
						Title:       "Prompt",
						Description: "Prompt describing the video to generate.",
					},
					"negativePrompt": {
						Type:        "string",
						Title:       "Negative Prompt",
						Description: "Aspects to discourage in the generated video.",
					},
					"cfgScale": {
						Type:        "number",
						Title:       "CFG Scale",
						Description: "Classifier-free guidance scale used during generation.",
						Minimum:     new(0.0),
						Maximum:     new(1.0),
						Default:     json.RawMessage(`0.5`),
					},
					"duration": {
						Type:        "number",
						Title:       "Duration",
						Description: "Duration of the generated video in seconds.",
						Default:     json.RawMessage(`5`),
					},
					"aspectRatio": {
						Type:        "string",
						Title:       "Aspect Ratio",
						Description: "Aspect ratio of the generated video.",
						Enum:        []any{"16:9", "9:16", "1:1"},
						Default:     json.RawMessage(`"16:9"`),
					},
					"cameraControlType": {
						Type:        "string",
						Title:       "Camera Control Type",
						Description: "Camera control preset to apply during generation.",
					},
					"cameraControlConfigHorizontal": {
						Type:        "number",
						Title:       "Camera Horizontal",
						Description: "Horizontal camera movement amount.",
					},
					"cameraControlConfigVertical": {
						Type:        "number",
						Title:       "Camera Vertical",
						Description: "Vertical camera movement amount.",
					},
					"cameraControlConfigPan": {
						Type:        "number",
						Title:       "Camera Pan",
						Description: "Camera pan movement amount.",
					},
					"cameraControlConfigTilt": {
						Type:        "number",
						Title:       "Camera Tilt",
						Description: "Camera tilt movement amount.",
					},
					"cameraControlConfigRoll": {
						Type:        "number",
						Title:       "Camera Roll",
						Description: "Camera roll movement amount.",
					},
					"cameraControlConfigZoom": {
						Type:        "number",
						Title:       "Camera Zoom",
						Description: "Camera zoom movement amount.",
					},
					"mode": {
						Type:        "string",
						Title:       "Mode",
						Description: "Generation mode (standard or pro).",
						Enum:        []any{"std", "pro"},
						Default:     json.RawMessage(`"std"`),
					},
					"version": {
						Type:        "string",
						Title:       "Version",
						Description: "PiAPI Kling model version. Version 2.1-master only supports pro mode.",
						Enum: []any{
							"2.6",
							"2.5",
							"2.1-master",
							"2.1",
							"1.6",
							"1.5",
						},
						Default: json.RawMessage(`"2.6"`),
					},
					"imageUrl": {
						Type:        "string",
						Title:       "Image URL",
						Description: "Optional reference image URL for image-to-video.",
						Format:      "uri",
					},
					"imageTailUrl": {
						Type:        "string",
						Title:       "Image Tail URL",
						Description: "Optional ending reference image URL.",
						Format:      "uri",
					},
					"motionBrushMaskUrl": {
						Type:        "string",
						Title:       "Motion Brush Mask URL",
						Description: "URL of the motion brush mask image.",
						Format:      "uri",
					},
					"motionBrushPointsStaticMasksX": {
						Type:        "number",
						Title:       "Motion Brush Static X",
						Description: "X coordinate for the static motion brush point.",
					},
					"motionBrushPointsStaticMasksY": {
						Type:        "number",
						Title:       "Motion Brush Static Y",
						Description: "Y coordinate for the static motion brush point.",
					},
					"motionBrushPointsDynamicMasksX": {
						Type:        "number",
						Title:       "Motion Brush Dynamic X",
						Description: "X coordinate for the dynamic motion brush point.",
					},
					"motionBrushPointsDynamicMasksY": {
						Type:        "number",
						Title:       "Motion Brush Dynamic Y",
						Description: "Y coordinate for the dynamic motion brush point.",
					},
					"effect": {
						Type:        "string",
						Title:       "Effect",
						Description: "Visual effect to apply to the generated video.",
						Enum:        []any{"squish", "expansion"},
						Default:     json.RawMessage(`"squish"`),
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
					"prompt",
					"aspectRatio",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"video": {
							Type:        "string",
							Description: "URL of the generated video.",
							Format:      "uri",
						},
						"image": {
							Type:        "string",
							Description: "URL of the generated cover image.",
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
