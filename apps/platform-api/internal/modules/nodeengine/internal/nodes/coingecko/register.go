package coingecko

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/coingecko/airdroptracker"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/coingecko/pricemonitor"
)

func Register() {
	nodeID := "coingecko"

	airdropTrackerNodeID := nodeID + "_airdrop_tracker"
	airdropTrackerNode := airdroptracker.NewCoingeckoAirdropTrackerNode(airdropTrackerNodeID)
	airdropTrackerValidator := jsonschema.New[airdroptracker.CoingeckoAirdropTrackerExecutorInput](airdropTrackerNode.GetInputSchema())
	airdropTrackerExecutor := airdroptracker.NewCoingeckoAirdropTrackerExecutor(airdropTrackerNodeID, airdropTrackerValidator)

	priceMonitorNodeID := nodeID + "_price_monitor"
	priceMonitorNode := pricemonitor.NewCoingeckoPriceMonitorNode(priceMonitorNodeID)
	priceMonitorValidator := jsonschema.New[pricemonitor.CoingeckoPriceMonitorExecutorInput](priceMonitorNode.GetInputSchema())
	priceMonitorExecutor := pricemonitor.NewCoingeckoPriceMonitorExecutor(priceMonitorNodeID, priceMonitorValidator)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "CoinGecko",
			Description: "Tools for tracking cryptocurrency markets via CoinGecko.",
			Icon: mcp.ServerIcon{
				Brand: "coingecko",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: airdropTrackerNode, Executor: airdropTrackerExecutor},
			{Node: priceMonitorNode, Executor: priceMonitorExecutor},
		},
	})
}
