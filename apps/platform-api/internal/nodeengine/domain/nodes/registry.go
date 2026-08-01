package nodes

import (
	"sync"
)

var (
	once     sync.Once
	instance *NodeRegistry
)

type NodeRegistry struct {
	nodesMap map[string]NodeManager
	nodes    []NodeManager
}

func GetRegistry() *NodeRegistry {
	once.Do(func() {
		instance = &NodeRegistry{
			nodesMap: make(map[string]NodeManager),
			nodes:    make([]NodeManager, 0),
		}
	})
	return instance
}

func (r *NodeRegistry) RegisterNode(node NodeManager) {
	if node.GetDisabled() {
		return
	}

	r.nodesMap[node.GetID()] = node
	r.nodes = append(r.nodes, node)
}

func (r *NodeRegistry) GetNode(id string) (NodeManager, bool) {
	node, ok := r.nodesMap[id]
	return node, ok
}

func (r *NodeRegistry) GetAllNodes() []NodeManager {
	return r.nodes
}

func RegisterNode(node NodeManager) {
	GetRegistry().RegisterNode(node)
}

func GetNode(id string) (NodeManager, bool) {
	return GetRegistry().GetNode(id)
}

func GetAllNodes() []NodeManager {
	return GetRegistry().GetAllNodes()
}
