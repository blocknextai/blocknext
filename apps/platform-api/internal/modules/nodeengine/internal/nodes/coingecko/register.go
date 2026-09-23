package coingecko

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/executors"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/nodes"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/coingecko/airdroptracker"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/coingecko/pricemonitor"
)

func Register() {
	nodeID := "coingecko"

	airdropTrackerNodeID := nodeID + "_airdrop_tracker"
	airdropTrackerNode := airdroptracker.NewCoingeckoAirdropTrackerNode(airdropTrackerNodeID)
	airdropTrackerValidator := jsonschema.New[airdroptracker.CoingeckoAirdropTrackerExecutorInput](airdropTrackerNode.GetInputSchema())
	airdropTrackerExecutor := airdroptracker.NewCoingeckoAirdropTrackerExecutor(airdropTrackerNodeID, airdropTrackerValidator)

	nodes.RegisterNode(airdropTrackerNode)
	executors.RegisterExecutor(airdropTrackerExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(airdropTrackerNode))

	priceMonitorNodeID := nodeID + "_price_monitor"
	priceMonitorNode := pricemonitor.NewCoingeckoPriceMonitorNode(priceMonitorNodeID)
	priceMonitorValidator := jsonschema.New[pricemonitor.CoingeckoPriceMonitorExecutorInput](priceMonitorNode.GetInputSchema())
	priceMonitorExecutor := pricemonitor.NewCoingeckoPriceMonitorExecutor(priceMonitorNodeID, priceMonitorValidator)

	nodes.RegisterNode(priceMonitorNode)
	executors.RegisterExecutor(priceMonitorExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(priceMonitorNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "CoinGecko",
		Description: "Tools for tracking cryptocurrency markets via CoinGecko.",
		Icon: mcp.ServerIcon{
			Brand: "coingecko",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			airdropTrackerNode,
			priceMonitorNode,
		},
	})
}
