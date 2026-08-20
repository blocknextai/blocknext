package createtrack

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type SoundCloudCreateTrackNode struct {
	nodes.Node
}

func NewSoundCloudCreateTrackNode(nodeID string) *SoundCloudCreateTrackNode {
	return &SoundCloudCreateTrackNode{
		ID:          nodeID,
		Kind:        nodes.NodeKindAction,
		Version:     "0.0.1",
		Name:        "SoundCloud Create Track",
		Description: "Upload a new track to SoundCloud.",
		Icon: nodes.NodeIcon{
			Brand: "soundcloud",
			Glyph: "music",
		},
		Inputs: []nodes.NodeHandle{
			{Key: "in"},
		},
		Outputs: []nodes.NodeHandle{
			{Key: "out"},
		},
		Categories:    []string{"Publishing"},
		SubCategories: []string{"SoundCloud"},
		Tags: []string{
			"audio",
			"upload",
			"create",
			"track",
			"music",
			"media",
			"publish",
			"share",
		},
		SupportedCredentials: []string{
			"soundcloud_oauth2",
		},
		InputSchema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"mp3Link": {
					Type:        "string",
					Title:       "MP3 Link",
					Description: "URL of the MP3 file to upload to SoundCloud.",
					Format:      "uri",
				},
				"title": {
					Type:        "string",
					Title:       "Title",
					Description: "Title of the SoundCloud track.",
				},
				"description": {
					Type:        "string",
					Title:       "Description",
					Description: "Optional description of the SoundCloud track.",
				},
			},
			Required: []string{
				"mp3Link",
				"title",
			},
		},
		OutputSchema: &gjs.Schema{
			Type: "array",
			Items: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"id": {
						Type:        "string",
						Description: "Identifier of the created SoundCloud track.",
					},
				},
			},
		},
		HasNaturalLanguage: true,
		Annotations: nodes.NodeAnnotations{
			Destructive: new(false),
		},
	}
}
