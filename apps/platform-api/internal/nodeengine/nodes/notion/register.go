package notion

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/notion/createpage"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/notion/getpage"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/notion/searchpages"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/notion/updatepage"
)

func Register() {
	nodeID := "notion"

	createPageNodeID := nodeID + "_create_page"
	createPageNode := createpage.NewNotionCreatePageNode(createPageNodeID)
	createPageValidator := jsonschema.New[createpage.NotionCreatePageExecutorInput](createPageNode.GetInputSchema())
	createPageExecutor := createpage.NewNotionCreatePageExecutor(createPageNodeID, createPageValidator)

	nodes.RegisterNode(createPageNode)
	executors.RegisterExecutor(createPageExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(createPageNode))

	updatePageNodeID := nodeID + "_update_page"
	updatePageNode := updatepage.NewNotionUpdatePageNode(updatePageNodeID)
	updatePageValidator := jsonschema.New[updatepage.NotionUpdatePageExecutorInput](updatePageNode.GetInputSchema())
	updatePageExecutor := updatepage.NewNotionUpdatePageExecutor(updatePageNodeID, updatePageValidator)

	nodes.RegisterNode(updatePageNode)
	executors.RegisterExecutor(updatePageExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(updatePageNode))

	searchPagesNodeID := nodeID + "_search_pages"
	searchPagesNode := searchpages.NewNotionSearchPagesNode(searchPagesNodeID)
	searchPagesValidator := jsonschema.New[searchpages.NotionSearchPagesExecutorInput](searchPagesNode.GetInputSchema())
	searchPagesExecutor := searchpages.NewNotionSearchPagesExecutor(searchPagesNodeID, searchPagesValidator)

	nodes.RegisterNode(searchPagesNode)
	executors.RegisterExecutor(searchPagesExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(searchPagesNode))

	getPageNodeID := nodeID + "_get_page"
	getPageNode := getpage.NewNotionGetPageNode(getPageNodeID)
	getPageValidator := jsonschema.New[getpage.NotionGetPageExecutorInput](getPageNode.GetInputSchema())
	getPageExecutor := getpage.NewNotionGetPageExecutor(getPageNodeID, getPageValidator)

	nodes.RegisterNode(getPageNode)
	executors.RegisterExecutor(getPageExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(getPageNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "Notion",
		Description: "Tools for managing Notion pages.",
		Icon: mcp.ServerIcon{
			Brand: "notion",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			createPageNode,
			updatePageNode,
			searchPagesNode,
			getPageNode,
		},
	})
}
