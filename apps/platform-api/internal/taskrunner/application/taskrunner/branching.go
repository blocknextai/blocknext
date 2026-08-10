package taskrunner

import (
	"strings"

	"github.com/blocknextai/go-packages/dag"
	nodeEngineDomainExecutors "github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	taskRunnerDomainTask "github.com/blocknextai/platform-api/internal/taskrunner/domain/task"
	"github.com/google/uuid"
)

const branchKeySeparator = "#"

func BuildBranchKey(nodeKey, handle string) string {
	var builder strings.Builder
	builder.WriteString(nodeKey)
	builder.WriteString(branchKeySeparator)
	builder.WriteString(handle)
	return builder.String()
}

type BranchView struct {
	TaskID    uuid.UUID
	Keys      map[string]string
	Absolute  []int
	InputKey  string
	collected map[string][]any
	matches   map[string][][]string
}

func (v *BranchView) Matches(pattern, template string, find func(string) [][]string) [][]string {
	if v == nil {
		return find(template)
	}
	if v.matches == nil {
		v.matches = make(map[string][][]string)
	}
	key := pattern + "\x00" + template
	if cached, ok := v.matches[key]; ok {
		return cached
	}
	found := find(template)
	v.matches[key] = found
	return found
}

func (v *BranchView) Collected(key string) ([]any, bool) {
	if v == nil || v.collected == nil {
		return nil, false
	}
	values, ok := v.collected[key]
	return values, ok
}

func (v *BranchView) Collect(key string, values []any) {
	if v == nil {
		return
	}
	if v.collected == nil {
		v.collected = make(map[string][]any)
	}
	v.collected[key] = values
}

func (v *BranchView) Input() string {
	if v == nil {
		return ""
	}
	return v.InputKey
}

func (v *BranchView) Task() uuid.UUID {
	if v == nil {
		return uuid.Nil
	}
	return v.TaskID
}

func (v *BranchView) StoreKey(nodeKey string) string {
	if v == nil {
		return nodeKey
	}
	if key, ok := v.Keys[nodeKey]; ok {
		return key
	}
	return nodeKey
}

func (v *BranchView) AbsoluteIndex(itemIndex int) int {
	if v == nil || v.Absolute == nil {
		return itemIndex
	}
	if itemIndex < 0 || itemIndex >= len(v.Absolute) {
		return itemIndex
	}
	return v.Absolute[itemIndex]
}

func (v *BranchView) ItemIndex(nodeKey string, itemIndex int) int {
	if v == nil {
		return itemIndex
	}
	if _, isBranched := v.Keys[nodeKey]; isBranched {
		return itemIndex
	}
	return v.AbsoluteIndex(itemIndex)
}

func edgeBetween(d *dag.DAG, parentID, childID string) *dag.Edge {
	for _, edge := range d.NodeEdges(parentID) {
		if edge.Target == childID {
			return &edge
		}
	}
	return nil
}

func (e *NodeExecutor) buildBranchView(task *taskRunnerDomainTask.Task, node *dag.Node) *BranchView {
	if task.DAG == nil {
		return nil
	}

	view := &BranchView{TaskID: task.ID, Keys: make(map[string]string)}
	var builder strings.Builder

	parents := task.DAG.NodeParents(node.ID)
	for _, parentID := range parents {
		parentNode := task.DAG.NodeByID(parentID)
		if parentNode == nil {
			continue
		}

		parentKey := BuildNodeKeyWithBuilder(&builder, parentNode.NodeID, parentNode.ID)
		storeKey := parentKey

		edge := edgeBetween(task.DAG, parentID, node.ID)
		if edge != nil && strings.TrimSpace(edge.SourceHandle) != "" {
			branchKey := BuildBranchKey(parentKey, edge.SourceHandle)
			if _, ok := e.outputStore.Get(task.ID, branchKey); ok {
				storeKey = branchKey
				view.Keys[parentKey] = branchKey
				if view.Absolute == nil {
					if absolute, ok := e.outputStore.Indexes(task.ID, branchKey); ok {
						view.Absolute = absolute
					}
				}
			}
		}

		if len(parents) == 1 {
			view.InputKey = storeKey
		}
	}

	return view
}

func branchReaches(store *OutputStore, taskID uuid.UUID, d *dag.DAG, parent *dag.Node, childID string) (bool, error) {
	edge := edgeBetween(d, parent.ID, childID)
	if edge == nil || strings.TrimSpace(edge.SourceHandle) == "" {
		return true, nil
	}

	executor, ok := nodeEngineDomainExecutors.GetExecutor(parent.NodeID)
	if !ok {
		return true, nil
	}
	if _, routes := executor.(nodeEngineDomainExecutors.BranchingExecutor); !routes {
		return true, nil
	}

	items, ok := store.Get(taskID, BuildBranchKey(BuildNodeKey(parent.NodeID, parent.ID), edge.SourceHandle))
	if !ok {
		return false, ErrBranchOutputsUnavailable
	}
	return len(items) > 0, nil
}

func routesItems(node *dag.Node) bool {
	executor, ok := nodeEngineDomainExecutors.GetExecutor(node.NodeID)
	if !ok {
		return false
	}
	_, routes := executor.(nodeEngineDomainExecutors.BranchingExecutor)
	return routes
}
