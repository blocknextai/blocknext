package notion

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/notion/createpage"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/notion/getpage"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/notion/searchpages"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/notion/updatepage"
)

func Register() {
	nodeID := "notion"

	createPageNodeID := nodeID + "_create_page"
	createPageNode := createpage.NewNotionCreatePageNode(createPageNodeID)
	createPageValidator := jsonschema.New[createpage.NotionCreatePageExecutorInput](createPageNode.GetInputSchema())
	createPageExecutor := createpage.NewNotionCreatePageExecutor(createPageNodeID, createPageValidator)

	updatePageNodeID := nodeID + "_update_page"
	updatePageNode := updatepage.NewNotionUpdatePageNode(updatePageNodeID)
	updatePageValidator := jsonschema.New[updatepage.NotionUpdatePageExecutorInput](updatePageNode.GetInputSchema())
	updatePageExecutor := updatepage.NewNotionUpdatePageExecutor(updatePageNodeID, updatePageValidator)

	searchPagesNodeID := nodeID + "_search_pages"
	searchPagesNode := searchpages.NewNotionSearchPagesNode(searchPagesNodeID)
	searchPagesValidator := jsonschema.New[searchpages.NotionSearchPagesExecutorInput](searchPagesNode.GetInputSchema())
	searchPagesExecutor := searchpages.NewNotionSearchPagesExecutor(searchPagesNodeID, searchPagesValidator)

	getPageNodeID := nodeID + "_get_page"
	getPageNode := getpage.NewNotionGetPageNode(getPageNodeID)
	getPageValidator := jsonschema.New[getpage.NotionGetPageExecutorInput](getPageNode.GetInputSchema())
	getPageExecutor := getpage.NewNotionGetPageExecutor(getPageNodeID, getPageValidator)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "Notion",
			Description: "Tools for managing Notion pages.",
			Icon: mcp.ServerIcon{
				Brand: "notion",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: createPageNode, Executor: createPageExecutor},
			{Node: updatePageNode, Executor: updatePageExecutor},
			{Node: searchPagesNode, Executor: searchPagesExecutor},
			{Node: getPageNode, Executor: getPageExecutor},
		},
	})
}
