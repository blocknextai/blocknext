package credentials

import (
	gjs "github.com/google/jsonschema-go/jsonschema"

	domain "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/credentials"
)

func NewCoingeckoAPICredential() *domain.Credential {
	return &domain.Credential{
		ID:          "coingecko_api",
		PlatformID:  "coingecko_api",
		Name:        "CoinGecko",
		Description: "CoinGecko API credentials for cryptocurrency market data.",
		Icon: domain.CredentialIcon{
			Brand: "coingecko",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"apiKey": {
					Type:        "string",
					Title:       "API Key",
					Description: "CoinGecko API key from your developer dashboard.",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"apiKey",
			},
		},
		SupportedNodes: &[]string{
			"coingecko_airdrop_tracker",
			"coingecko_price_monitor",
		},
	}
}
