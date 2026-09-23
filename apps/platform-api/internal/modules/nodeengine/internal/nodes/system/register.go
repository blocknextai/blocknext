package system

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/executors"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/nodes"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/system/annotation"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/system/condition"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/system/sleep"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/system/starter"
)

func Register() {
	nodeID := "system"

	conditionNodeID := nodeID + "_condition"
	conditionNode := condition.NewConditionNode(conditionNodeID)
	conditionValidator := jsonschema.New[condition.ConditionExecutorInput](conditionNode.GetInputSchema())
	conditionExecutor := condition.NewConditionExecutor(conditionNodeID, conditionValidator)

	nodes.RegisterNode(conditionNode)
	executors.RegisterExecutor(conditionExecutor)

	sleepNodeID := nodeID + "_sleep"
	sleepNode := sleep.NewSleepNode(sleepNodeID)
	sleepValidator := jsonschema.New[sleep.SleepExecutorInput](sleepNode.GetInputSchema())
	sleepExecutor := sleep.NewSleepExecutor(sleepNodeID, sleepValidator)

	nodes.RegisterNode(sleepNode)
	executors.RegisterExecutor(sleepExecutor)

	starterNodeID := nodeID + "_starter"
	starterNode := starter.NewStarterNode(starterNodeID)
	starterExecutor := starter.NewStarterExecutor(starterNodeID)

	nodes.RegisterNode(starterNode)
	executors.RegisterExecutor(starterExecutor)

	annotationNodeID := nodeID + "_annotation"
	annotationNode := annotation.NewAnnotationNode(annotationNodeID)

	nodes.RegisterNode(annotationNode)

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "System",
		Description: "System-level workflow tools.",
		Version:     "0.0.1",
		Tools: []nodes.NodeManager{
			conditionNode,
			sleepNode,
			starterNode,
		},
	})
}
