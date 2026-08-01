package pricemonitor

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type CoingeckoPriceMonitorNode struct {
	nodes.Node
}

func NewCoingeckoPriceMonitorNode(nodeID string) *CoingeckoPriceMonitorNode {
	return &CoingeckoPriceMonitorNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Coingecko Price Monitor",
			Description: "Monitor cryptocurrency prices via CoinGecko.",
			Icon: nodes.NodeIcon{
				Light: nodeID,
				Dark:  nodeID,
			},
			Categories:    []string{"Blockchain"},
			SubCategories: []string{"CoinGecko"},
			Tags: []string{
				"price",
				"crypto",
				"monitoring",
				"blockchain",
				"alert",
				"tracking",
				"web3",
			},
			SupportedCredentials: []string{
				"coingecko_api",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"coinId": {
						Type:        "string",
						Title:       "Coin ID",
						Description: "CoinGecko coin identifier (e.g. bitcoin, ethereum).",
						Default:     json.RawMessage(`"bitcoin"`),
					},
					"currency": {
						Type:        "string",
						Title:       "Currency",
						Description: "Quote currency code (e.g. usd, eur).",
						Default:     json.RawMessage(`"usd"`),
					},
					"alertType": {
						Type:        "string",
						Title:       "Alert Type",
						Description: "Type of alert: price_above, price_below, change_above, change_below, or info.",
						Default:     json.RawMessage(`"info"`),
					},
					"thresholdValue": {
						Type:        "number",
						Title:       "Threshold Value",
						Description: "Threshold value used together with the alert type.",
						Default:     json.RawMessage(`0`),
					},
				},
				Required: []string{
					"coinId",
					"currency",
					"alertType",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"coinId": {
							Type:        "string",
							Description: "CoinGecko coin identifier that was queried.",
						},
						"currency": {
							Type:        "string",
							Description: "Quote currency that was used.",
						},
						"currentPrice": {
							Type:        "number",
							Description: "Current price in the requested currency.",
						},
						"change24h": {
							Type:        "number",
							Description: "Percentage price change in the last 24 hours.",
						},
						"volume24h": {
							Type:        "number",
							Description: "Trading volume in the last 24 hours.",
						},
						"marketCap": {
							Type:        "number",
							Description: "Market capitalization in the requested currency.",
						},
						"alertType": {
							Type:        "string",
							Description: "Alert type that was evaluated.",
						},
						"thresholdValue": {
							Type:        "number",
							Description: "Threshold value that was evaluated.",
						},
						"shouldAlert": {
							Type:        "boolean",
							Description: "Whether the alert condition was triggered.",
						},
						"changeDirection": {
							Type:        "string",
							Description: "Direction of the price change: up, down, or neutral.",
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
