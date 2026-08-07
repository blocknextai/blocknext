package airdroptracker

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type CoingeckoAirdropTrackerNode struct {
	nodes.Node
}

func NewCoingeckoAirdropTrackerNode(nodeID string) *CoingeckoAirdropTrackerNode {
	return &CoingeckoAirdropTrackerNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Coingecko Airdrop Tracker",
			Description: "Track cryptocurrency airdrops via CoinGecko.",
			Icon: nodes.NodeIcon{
				Brand: "coingecko",
				Glyph: "target",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"Blockchain"},
			SubCategories: []string{"CoinGecko"},
			Tags: []string{
				"airdrop",
				"crypto",
				"tracking",
				"blockchain",
				"monitoring",
				"alert",
				"web3",
			},
			SupportedCredentials: []string{
				"coingecko_api",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"projectName": {
						Type:        "string",
						Title:       "Project Name",
						Description: "Filter airdrops by project name (case-insensitive substring match).",
					},
					"statusFilter": {
						Type:        "string",
						Title:       "Status Filter",
						Description: "Airdrop status to filter by (e.g. ongoing, upcoming, ended).",
						Default:     json.RawMessage(`"ongoing"`),
					},
					"minReward": {
						Type:        "number",
						Title:       "Minimum Reward",
						Description: "Minimum reward value required to include the airdrop.",
						Default:     json.RawMessage(`0`),
					},
				},
				Required: []string{
					"projectName",
					"statusFilter",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"airdrops": {
							Type: "array",
							Items: &gjs.Schema{
								Type: "object",
								Properties: map[string]*gjs.Schema{
									"id":              {Type: "string"},
									"name":            {Type: "string"},
									"description":     {Type: "string"},
									"status":          {Type: "string"},
									"startDate":       {Type: "string"},
									"endDate":         {Type: "string"},
									"projectName":     {Type: "string"},
									"projectWebsite":  {Type: "string"},
									"projectTwitter":  {Type: "string"},
									"projectTelegram": {Type: "string"},
									"projectDiscord":  {Type: "string"},
									"rewardAmount":    {Type: "string"},
									"rewardCurrency":  {Type: "string"},
									"rewardValue":     {Type: "number"},
									"requirements": {
										Type:  "array",
										Items: &gjs.Schema{Type: "string"},
									},
									"instructions": {Type: "string"},
								},
							},
						},
						"totalCount": {
							Type:        "integer",
							Description: "Total number of airdrops returned after filtering.",
						},
						"projectFilter": {
							Type:        "string",
							Description: "Project filter applied to the request.",
						},
						"statusFilter": {
							Type:        "string",
							Description: "Status filter applied to the request.",
						},
						"minReward": {
							Type:        "number",
							Description: "Minimum reward filter applied to the request.",
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
