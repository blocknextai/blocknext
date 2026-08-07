// Package dag builds and traverses directed acyclic graphs of workflow nodes
// and edges. It provides topological ordering with type-based priorities,
// reachability analysis, orphan detection, and conditional branching.
package dag

import (
	"container/heap"
	"log/slog"
	"strings"
)

// DAG represents a directed acyclic graph of nodes and edges, precomputed for
// reachability, topological ordering, and conditional traversal.
type DAG struct {
	adjacencyList    map[string][]string
	reverseAdjacency map[string][]string
	nodeTypes        map[string]string
	nodeData         map[string]*Node
	edgeData         map[string][]Edge
	startNodes       []string
	sortedNodeIDs    []string
	sortedNodes      []Node
	sortedEdges      []Edge
	reachableCache   map[string]bool
}

// New builds a DAG from the given nodes and edges, validating them and
// computing the topological execution order. It returns an error if the input
// is empty or invalid, or if a cycle is detected.
func New(nodes []Node, edges []Edge) (*DAG, error) {
	if len(nodes) == 0 {
		return nil, ErrEmptyDAG
	}

	dag := &DAG{
		adjacencyList:    make(map[string][]string),
		reverseAdjacency: make(map[string][]string),
		nodeTypes:        make(map[string]string),
		nodeData:         make(map[string]*Node),
		edgeData:         make(map[string][]Edge),
		startNodes:       make([]string, 0),
		sortedNodeIDs:    make([]string, 0),
		sortedNodes:      make([]Node, 0),
		sortedEdges:      make([]Edge, 0),
	}

	for i := range nodes {
		node := &nodes[i]
		if err := dag.addNode(node); err != nil {
			return nil, err
		}
	}

	for _, edge := range edges {
		if err := dag.addEdge(edge); err != nil {
			return nil, err
		}
	}

	dag.findStartNodes()

	if err := dag.computeSortedOrder(); err != nil {
		return nil, err
	}

	return dag, nil
}

func (d *DAG) addNode(node *Node) error {
	if node == nil || strings.TrimSpace(node.ID) == "" {
		return ErrInvalidNode
	}

	d.adjacencyList[node.ID] = make([]string, 0)
	d.reverseAdjacency[node.ID] = make([]string, 0)
	d.nodeTypes[node.ID] = node.NodeID
	d.nodeData[node.ID] = node
	d.edgeData[node.ID] = make([]Edge, 0)

	return nil
}

func (d *DAG) addEdge(edge Edge) error {
	if strings.TrimSpace(edge.Source) == "" || strings.TrimSpace(edge.Target) == "" {
		return ErrInvalidNode
	}

	orphanedEdge := false
	if _, exists := d.nodeData[edge.Source]; !exists {
		slog.Warn("orphaned edge detected - skipping",
			"package", "dag",
			"source_node", edge.Source,
			"target_node", edge.Target,
			"reason", "source node does not exist",
		)
		orphanedEdge = true
	}
	if _, exists := d.nodeData[edge.Target]; !exists {
		slog.Warn("orphaned edge detected - skipping",
			"package", "dag",
			"source_node", edge.Source,
			"target_node", edge.Target,
			"reason", "target node does not exist",
		)
		orphanedEdge = true
	}

	if orphanedEdge {
		return nil
	}

	d.adjacencyList[edge.Source] = append(d.adjacencyList[edge.Source], edge.Target)
	d.reverseAdjacency[edge.Target] = append(d.reverseAdjacency[edge.Target], edge.Source)
	d.edgeData[edge.Source] = append(d.edgeData[edge.Source], edge)

	return nil
}

func (d *DAG) findStartNodes() {
	incomingCount := make(map[string]int)

	for nodeID := range d.nodeData {
		incomingCount[nodeID] = 0
	}

	for _, targets := range d.adjacencyList {
		for _, target := range targets {
			incomingCount[target]++
		}
	}

	d.startNodes = make([]string, 0)
	for nodeID, count := range incomingCount {
		if count == 0 && d.nodeTypes[nodeID] == "system_starter" {
			d.startNodes = append(d.startNodes, nodeID)
		}
	}
}

func (d *DAG) findReachableNodes() map[string]bool {
	if d.reachableCache != nil {
		return d.reachableCache
	}

	reachable := make(map[string]bool)

	for _, startNode := range d.startNodes {
		d.dfsReachable(startNode, reachable)
	}

	d.reachableCache = reachable

	return reachable
}

func (d *DAG) dfsReachable(nodeID string, reachable map[string]bool) {
	if reachable[nodeID] {
		return
	}

	reachable[nodeID] = true

	for _, child := range d.adjacencyList[nodeID] {
		d.dfsReachable(child, reachable)
	}
}

// OrphanedNodes returns the IDs of nodes that are not reachable from any start node.
func (d *DAG) OrphanedNodes() []string {
	reachable := d.findReachableNodes()
	var orphaned []string

	for nodeID := range d.nodeData {
		if !reachable[nodeID] {
			orphaned = append(orphaned, nodeID)
		}
	}

	return orphaned
}

// HasOrphanedNodes reports whether the DAG contains any unreachable nodes.
func (d *DAG) HasOrphanedNodes() bool {
	return len(d.OrphanedNodes()) > 0
}

// ReachableNodes returns the nodes reachable from the start nodes.
func (d *DAG) ReachableNodes() []Node {
	reachable := d.findReachableNodes()
	var nodes []Node

	for nodeID := range reachable {
		if node, exists := d.nodeData[nodeID]; exists {
			nodes = append(nodes, *node)
		}
	}

	return nodes
}

// ReachableNodeIDs returns the IDs of nodes reachable from the start nodes.
func (d *DAG) ReachableNodeIDs() []string {
	reachable := d.findReachableNodes()
	var nodeIDs []string

	for nodeID := range reachable {
		nodeIDs = append(nodeIDs, nodeID)
	}

	return nodeIDs
}

// IsNodeReachable reports whether the node with the given ID is reachable from a start node.
func (d *DAG) IsNodeReachable(nodeID string) bool {
	reachable := d.findReachableNodes()
	return reachable[nodeID]
}

func (d *DAG) computeSortedOrder() error {
	if len(d.nodeData) == 0 {
		return ErrEmptyDAG
	}

	reachable := d.findReachableNodes()

	orphanedNodes := d.OrphanedNodes()
	if len(orphanedNodes) > 0 {
		slog.Warn("orphaned nodes detected - will be excluded from execution",
			"package", "dag",
			"orphaned_nodes", orphanedNodes,
			"count", len(orphanedNodes),
		)
	}

	if len(reachable) == 0 {
		return ErrOrphanedNode
	}

	inDegree := make(map[string]int)

	for nodeID := range reachable {
		inDegree[nodeID] = 0
	}

	for _, neighbors := range d.adjacencyList {
		for _, neighbor := range neighbors {
			if reachable[neighbor] {
				inDegree[neighbor]++
			}
		}
	}

	pq := &PriorityQueue{}
	heap.Init(pq)

	for nodeID, degree := range inDegree {
		if degree == 0 {
			heap.Push(pq, &Item{
				ID:       nodeID,
				Priority: d.getNodePriority(nodeID),
			})
		}
	}

	var result []string

	for pq.Len() > 0 {
		item := heap.Pop(pq).(*Item)
		nodeID := item.ID
		result = append(result, nodeID)

		for _, neighbor := range d.adjacencyList[nodeID] {
			if reachable[neighbor] {
				inDegree[neighbor]--
				if inDegree[neighbor] == 0 {
					heap.Push(pq, &Item{
						ID:       neighbor,
						Priority: d.getNodePriority(neighbor),
					})
				}
			}
		}
	}

	if len(result) != len(reachable) {
		return ErrCycleDetected
	}

	d.sortedNodeIDs = result
	d.sortedNodes = make([]Node, len(result))
	d.sortedEdges = make([]Edge, 0)

	for i, nodeID := range result {
		node := d.nodeData[nodeID]
		d.sortedNodes[i] = *node
		d.sortedEdges = append(d.sortedEdges, d.edgeData[nodeID]...)
	}

	return nil
}

func (d *DAG) getNodePriority(nodeID string) int {
	nodeType := d.nodeTypes[nodeID]
	switch nodeType {
	case "system_starter":
		return 0
	case "system_condition":
		return 1
	case "system_action":
		return 2
	default:
		return 3
	}
}

// StartNodes returns the IDs of the graph's start nodes.
func (d *DAG) StartNodes() []string {
	return d.startNodes
}

// NodeByID returns the node with the given ID, or nil if it does not exist.
func (d *DAG) NodeByID(nodeID string) *Node {
	return d.nodeData[nodeID]
}

// NodeChildren returns the IDs of the direct children of the node with the given ID.
func (d *DAG) NodeChildren(nodeID string) []string {
	return d.adjacencyList[nodeID]
}

// NodeEdges returns the outgoing edges of the node with the given ID.
func (d *DAG) NodeEdges(nodeID string) []Edge {
	return d.edgeData[nodeID]
}

// NodeCount returns the total number of nodes in the DAG.
func (d *DAG) NodeCount() int {
	return len(d.nodeData)
}

// EdgeCount returns the total number of edges in the DAG.
func (d *DAG) EdgeCount() int {
	count := 0
	for _, edges := range d.edgeData {
		count += len(edges)
	}
	return count
}

// Nodes returns the reachable nodes in topological execution order.
func (d *DAG) Nodes() []Node {
	return d.sortedNodes
}

// Edges returns the edges of the reachable nodes in topological execution order.
func (d *DAG) Edges() []Edge {
	return d.sortedEdges
}

// Validate reports an error if the DAG has no nodes in its sorted order.
func (d *DAG) Validate() error {
	if len(d.sortedNodeIDs) == 0 {
		return ErrEmptyDAG
	}
	return nil
}

// IncomingDegree returns the number of edges pointing to the node with the given ID.
func (d *DAG) IncomingDegree(nodeID string) int {
	if parents, exists := d.reverseAdjacency[nodeID]; exists {
		return len(parents)
	}
	return 0
}

// NodeParents returns the IDs of the direct parents of the node with the given ID.
func (d *DAG) NodeParents(nodeID string) []string {
	if parents, exists := d.reverseAdjacency[nodeID]; exists {
		return parents
	}
	return []string{}
}

// AllNodeIDs returns the IDs of every node in the DAG, including orphans.
func (d *DAG) AllNodeIDs() []string {
	var nodeIDs []string
	for nodeID := range d.nodeData {
		nodeIDs = append(nodeIDs, nodeID)
	}
	return nodeIDs
}
