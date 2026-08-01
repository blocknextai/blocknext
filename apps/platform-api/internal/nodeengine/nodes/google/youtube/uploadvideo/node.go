package uploadvideo

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type YouTubeUploadVideoNode struct {
	nodes.Node
}

func NewYouTubeUploadVideoNode(nodeID string) *YouTubeUploadVideoNode {
	return &YouTubeUploadVideoNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "YouTube Upload Video",
			Description: "Upload a video file to YouTube with the given metadata and privacy setting.",
			Icon: nodes.NodeIcon{
				Light: nodeID,
				Dark:  nodeID,
			},
			Categories:    []string{"Publishing"},
			SubCategories: []string{"YouTube"},
			Tags: []string{
				"google",
				"video",
				"upload",
				"publish",
				"media",
				"share",
			},
			SupportedCredentials: []string{
				"youtube_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"title": {
						Type:        "string",
						Title:       "Title",
						Description: "Title of the YouTube video.",
						Default:     json.RawMessage(`"via BlockNext"`),
					},
					"description": {
						Type:        "string",
						Title:       "Description",
						Description: "Description shown below the video on YouTube.",
					},
					"categoryId": {
						Type:        "string",
						Title:       "Category",
						Description: "YouTube video category identifier.",
						Enum:        []any{"1", "2", "10", "15", "17", "19", "20", "22", "23", "24", "25", "26", "27", "28", "29"},
						Default:     json.RawMessage(`"22"`),
					},
					"privacy": {
						Type:        "string",
						Title:       "Privacy",
						Description: "Privacy status for the uploaded video.",
						Enum:        []any{"public", "private", "unlisted"},
						Default:     json.RawMessage(`"public"`),
					},
					"videoUrl": {
						Type:        "string",
						Title:       "Video URL",
						Description: "URL of the source video file to upload.",
						Format:      "uri",
					},
				},
				Required: []string{
					"title",
					"categoryId",
					"privacy",
					"videoUrl",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"status": {
							Type:        "boolean",
							Description: "Whether the upload completed successfully.",
						},
						"videoId": {
							Type:        "string",
							Description: "Identifier of the uploaded YouTube video.",
						},
						"videoUrl": {
							Type:        "string",
							Format:      "uri",
							Description: "Public YouTube URL of the uploaded video.",
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
