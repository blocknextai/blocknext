package createplaylist

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type SoundCloudCreatePlaylistNode struct {
	nodes.Node
}

func NewSoundCloudCreatePlaylistNode(nodeID string) *SoundCloudCreatePlaylistNode {
	return &SoundCloudCreatePlaylistNode{
		ID:          nodeID,
		Kind:        nodes.NodeKindAction,
		Version:     "0.0.1",
		Name:        "SoundCloud Create Playlist",
		Description: "Create a new playlist on SoundCloud.",
		Icon: nodes.NodeIcon{
			Brand: "soundcloud",
			Glyph: "list",
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
			"playlist",
			"create",
			"collection",
			"music",
			"audio",
			"media",
			"share",
		},
		SupportedCredentials: []string{
			"soundcloud_oauth2",
		},
		InputSchema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"title": {
					Type:        "string",
					Title:       "Title",
					Description: "Title of the SoundCloud playlist.",
				},
				"description": {
					Type:        "string",
					Title:       "Description",
					Description: "Optional description of the SoundCloud playlist.",
				},
				"trackIds": {
					Type:        "string",
					Title:       "Track IDs",
					Description: "Identifiers of the tracks to include in the playlist.",
				},
			},
			Required: []string{
				"title",
			},
		},
		OutputSchema: &gjs.Schema{
			Type: "array",
			Items: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"status": {
						Type:        "boolean",
						Description: "Whether the SoundCloud playlist was created successfully.",
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
