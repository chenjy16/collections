package graph

import (
	"fmt"
	"strings"

	"github.com/chenjianyu/collections/container/common"
	"github.com/chenjianyu/collections/container/set"
)

// MutableValueGraph is a mutable implementation of ValueGraph
type MutableValueGraph[N comparable, V any] struct {
	directed       bool
	allowSelfLoops bool
	nodeOrder      ElementOrder
	nodes          set.Set[N]
	// edgeValues maps edge endpoints to their values
	edgeValues map[EndpointPair[N]]V
	// adjacencyMap maps each node to its adjacent nodes
	adjacencyMap map[N]set.Set[N]
	// For directed graphs, we also maintain predecessors
	predecessorMap map[N]set.Set[N]
}

// NewMutableValueGraph creates a new mutable value graph
func NewMutableValueGraph[N comparable, V any](directed, allowSelfLoops bool, nodeOrder ElementOrder) *MutableValueGraph[N, V] {
	return &MutableValueGraph[N, V]{
		directed:       directed,
		allowSelfLoops: allowSelfLoops,
		nodeOrder:      nodeOrder,
		nodes:          set.New[N](),
		edgeValues:     make(map[EndpointPair[N]]V),
		adjacencyMap:   make(map[N]set.Set[N]),
		predecessorMap: make(map[N]set.Set[N]),
	}
}

// Nodes returns all nodes in this graph
func (g *MutableValueGraph[N, V]) Nodes() set.Set[N] {
	result := set.New[N]()
	g.nodes.ForEach(func(node N) {
		result.Add(node)
	})
	return result
}

// Edges returns all edges in this graph as endpoint pairs
func (g *MutableValueGraph[N, V]) Edges() set.Set[EndpointPair[N]] {
	result := set.New[EndpointPair[N]]()
	
	for edge := range g.edgeValues {
		result.Add(edge)
	}
	
	return result
}

// IsDirected returns true if this is a directed graph
func (g *MutableValueGraph[N, V]) IsDirected() bool {
	return g.directed
}

// AllowsSelfLoops returns true if this graph allows self-loops
func (g *MutableValueGraph[N, V]) AllowsSelfLoops() bool {
	return g.allowSelfLoops
}

// NodeOrder returns the ordering of nodes in this graph
func (g *MutableValueGraph[N, V]) NodeOrder() ElementOrder {
	return g.nodeOrder
}

// AddNode adds a node to this graph
func (g *MutableValueGraph[N, V]) AddNode(node N) bool {
	if g.nodes.Contains(node) {
		return false
	}
	
	g.nodes.Add(node)
	g.adjacencyMap[node] = set.New[N]()
	if g.directed {
		g.predecessorMap[node] = set.New[N]()
	}
	return true
}

// PutEdge adds an edge between two nodes
func (g *MutableValueGraph[N, V]) PutEdge(nodeU, nodeV N) bool {
	var zeroValue V
	_, existed := g.PutEdgeValue(nodeU, nodeV, zeroValue)
	return !existed
}

// RemoveNode removes a node and all its incident edges
func (g *MutableValueGraph[N, V]) RemoveNode(node N) bool {
	if !g.nodes.Contains(node) {
		return false
	}
	
	// Remove all edges incident to this node
	successors := g.adjacencyMap[node]
	successors.ForEach(func(successor N) {
		g.RemoveEdge(node, successor)
	})
	
	if g.directed {
		predecessors := g.predecessorMap[node]
		predecessors.ForEach(func(predecessor N) {
			g.RemoveEdge(predecessor, node)
		})
	}
	
	// Remove the node itself
	g.nodes.Remove(node)
	delete(g.adjacencyMap, node)
	if g.directed {
		delete(g.predecessorMap, node)
	}
	
	return true
}

// RemoveEdge removes an edge between two nodes
func (g *MutableValueGraph[N, V]) RemoveEdge(nodeU, nodeV N) bool {
	if !g.nodes.Contains(nodeU) || !g.nodes.Contains(nodeV) {
		return false
	}
	
	edge := NewEndpointPair(nodeU, nodeV)
	if _, exists := g.edgeValues[edge]; !exists {
		if !g.directed {
			// Try the reverse edge for undirected graphs
			reverseEdge := NewEndpointPair(nodeV, nodeU)
			if _, exists := g.edgeValues[reverseEdge]; !exists {
				return false
			}
			edge = reverseEdge
		} else {
			return false
		}
	}
	
	delete(g.edgeValues, edge)
	g.adjacencyMap[nodeU].Remove(nodeV)
	if g.directed {
		g.predecessorMap[nodeV].Remove(nodeU)
	} else {
		g.adjacencyMap[nodeV].Remove(nodeU)
	}
	
	return true
}

// Successors returns the successors of a node
func (g *MutableValueGraph[N, V]) Successors(node N) set.Set[N] {
	if !g.nodes.Contains(node) {
		panic(fmt.Sprintf("Node %v is not in the graph", node))
	}
	
	result := set.New[N]()
	g.adjacencyMap[node].ForEach(func(successor N) {
		result.Add(successor)
	})
	return result
}

// Predecessors returns the predecessors of a node
func (g *MutableValueGraph[N, V]) Predecessors(node N) set.Set[N] {
	if !g.nodes.Contains(node) {
		panic(fmt.Sprintf("Node %v is not in the graph", node))
	}
	
	result := set.New[N]()
	if g.directed {
		g.predecessorMap[node].ForEach(func(predecessor N) {
			result.Add(predecessor)
		})
	} else {
		// For undirected graphs, predecessors are the same as successors
		g.adjacencyMap[node].ForEach(func(adjacent N) {
			result.Add(adjacent)
		})
	}
	return result
}

// AdjacentNodes returns all nodes adjacent to the given node
func (g *MutableValueGraph[N, V]) AdjacentNodes(node N) set.Set[N] {
	if !g.nodes.Contains(node) {
		panic(fmt.Sprintf("Node %v is not in the graph", node))
	}
	
	result := set.New[N]()
	
	// Add successors
	g.adjacencyMap[node].ForEach(func(successor N) {
		result.Add(successor)
	})
	
	// Add predecessors (for directed graphs, they might be different)
	if g.directed {
		g.predecessorMap[node].ForEach(func(predecessor N) {
			result.Add(predecessor)
		})
	}
	
	return result
}

// IncidentEdges returns all edges incident to the given node
func (g *MutableValueGraph[N, V]) IncidentEdges(node N) set.Set[EndpointPair[N]] {
	if !g.nodes.Contains(node) {
		panic(fmt.Sprintf("Node %v is not in the graph", node))
	}
	
	result := set.New[EndpointPair[N]]()
	
	for edge := range g.edgeValues {
		if edge.NodeU == node || edge.NodeV == node {
			result.Add(edge)
		}
	}
	
	return result
}

// Degree returns the degree of a node
func (g *MutableValueGraph[N, V]) Degree(node N) int {
	if !g.nodes.Contains(node) {
		panic(fmt.Sprintf("Node %v is not in the graph", node))
	}
	
	if g.directed {
		return g.InDegree(node) + g.OutDegree(node)
	}
	return g.adjacencyMap[node].Size()
}

// InDegree returns the in-degree of a node
func (g *MutableValueGraph[N, V]) InDegree(node N) int {
	if !g.nodes.Contains(node) {
		panic(fmt.Sprintf("Node %v is not in the graph", node))
	}
	
	if g.directed {
		return g.predecessorMap[node].Size()
	}
	return g.adjacencyMap[node].Size()
}

// OutDegree returns the out-degree of a node
func (g *MutableValueGraph[N, V]) OutDegree(node N) int {
	if !g.nodes.Contains(node) {
		panic(fmt.Sprintf("Node %v is not in the graph", node))
	}
	
	return g.adjacencyMap[node].Size()
}

// HasEdgeConnecting returns true if there's an edge between two nodes
func (g *MutableValueGraph[N, V]) HasEdgeConnecting(nodeU, nodeV N) bool {
	if !g.nodes.Contains(nodeU) || !g.nodes.Contains(nodeV) {
		return false
	}
	
	edge := NewEndpointPair(nodeU, nodeV)
	if _, exists := g.edgeValues[edge]; exists {
		return true
	}
	
	if !g.directed {
		// Check reverse edge for undirected graphs
		reverseEdge := NewEndpointPair(nodeV, nodeU)
		_, exists := g.edgeValues[reverseEdge]
		return exists
	}
	
	return false
}

// EdgeValue returns the value associated with an edge
func (g *MutableValueGraph[N, V]) EdgeValue(nodeU, nodeV N) (V, bool) {
	var zeroValue V
	
	if !g.nodes.Contains(nodeU) || !g.nodes.Contains(nodeV) {
		return zeroValue, false
	}
	
	edge := NewEndpointPair(nodeU, nodeV)
	if value, exists := g.edgeValues[edge]; exists {
		return value, true
	}
	
	if !g.directed {
		// Check reverse edge for undirected graphs
		reverseEdge := NewEndpointPair(nodeV, nodeU)
		if value, exists := g.edgeValues[reverseEdge]; exists {
			return value, true
		}
	}
	
	return zeroValue, false
}

// EdgeValueOrDefault returns the value associated with an edge or a default value
func (g *MutableValueGraph[N, V]) EdgeValueOrDefault(nodeU, nodeV N, defaultValue V) V {
	if value, exists := g.EdgeValue(nodeU, nodeV); exists {
		return value
	}
	return defaultValue
}

// PutEdgeValue adds an edge with a value
func (g *MutableValueGraph[N, V]) PutEdgeValue(nodeU, nodeV N, value V) (V, bool) {
	var zeroValue V
	
	// Check self-loop constraint
	if !g.allowSelfLoops && nodeU == nodeV {
		panic("Self-loops are not allowed in this graph")
	}
	
	// Add nodes if they don't exist
	g.AddNode(nodeU)
	g.AddNode(nodeV)
	
	edge := NewEndpointPair(nodeU, nodeV)
	
	// Check if edge already exists
	if existingValue, exists := g.edgeValues[edge]; exists {
		g.edgeValues[edge] = value
		return existingValue, true
	}
	
	// For undirected graphs, check if reverse edge exists
	if !g.directed {
		reverseEdge := NewEndpointPair(nodeV, nodeU)
		if existingValue, exists := g.edgeValues[reverseEdge]; exists {
			g.edgeValues[reverseEdge] = value
			return existingValue, true
		}
	}
	
	// Add new edge
	g.edgeValues[edge] = value
	g.adjacencyMap[nodeU].Add(nodeV)
	if g.directed {
		g.predecessorMap[nodeV].Add(nodeU)
	} else {
		g.adjacencyMap[nodeV].Add(nodeU)
	}
	
	return zeroValue, false
}

// AsGraph returns a view of this value graph as a basic graph
func (g *MutableValueGraph[N, V]) AsGraph() Graph[N] {
	return &valueGraphAsGraph[N, V]{g}
}

// Size returns the number of nodes in the graph
func (g *MutableValueGraph[N, V]) Size() int {
	return g.nodes.Size()
}

// IsEmpty returns true if the graph has no nodes
func (g *MutableValueGraph[N, V]) IsEmpty() bool {
	return g.nodes.IsEmpty()
}

// Clear removes all nodes and edges from the graph
func (g *MutableValueGraph[N, V]) Clear() {
	g.nodes.Clear()
	g.edgeValues = make(map[EndpointPair[N]]V)
	g.adjacencyMap = make(map[N]set.Set[N])
	g.predecessorMap = make(map[N]set.Set[N])
}

// Contains returns true if the graph contains the given node
func (g *MutableValueGraph[N, V]) Contains(node N) bool {
	return g.nodes.Contains(node)
}

// ToSlice returns all nodes as a slice
func (g *MutableValueGraph[N, V]) ToSlice() []N {
	return g.nodes.ToSlice()
}

// Iterator returns an iterator over the nodes
func (g *MutableValueGraph[N, V]) Iterator() common.Iterator[N] {
	return g.nodes.Iterator()
}

// ForEach applies a function to each node
func (g *MutableValueGraph[N, V]) ForEach(fn func(N)) {
	g.nodes.ForEach(fn)
}

// String returns a string representation of the graph
func (g *MutableValueGraph[N, V]) String() string {
	var sb strings.Builder
	sb.WriteString("ValueGraph{")
	
	if g.directed {
		sb.WriteString("directed, ")
	} else {
		sb.WriteString("undirected, ")
	}
	
	sb.WriteString(fmt.Sprintf("nodes=%d, edges=%d", g.Size(), len(g.edgeValues)))
	sb.WriteString("}")
	
	return sb.String()
}

// valueGraphAsGraph is an adapter that presents a ValueGraph as a Graph
type valueGraphAsGraph[N comparable, V any] struct {
	valueGraph *MutableValueGraph[N, V]
}

func (g *valueGraphAsGraph[N, V]) Nodes() set.Set[N] {
	return g.valueGraph.Nodes()
}

func (g *valueGraphAsGraph[N, V]) Edges() set.Set[EndpointPair[N]] {
	return g.valueGraph.Edges()
}

func (g *valueGraphAsGraph[N, V]) IsDirected() bool {
	return g.valueGraph.IsDirected()
}

func (g *valueGraphAsGraph[N, V]) AllowsSelfLoops() bool {
	return g.valueGraph.AllowsSelfLoops()
}

func (g *valueGraphAsGraph[N, V]) NodeOrder() ElementOrder {
	return g.valueGraph.NodeOrder()
}

func (g *valueGraphAsGraph[N, V]) AddNode(node N) bool {
	return g.valueGraph.AddNode(node)
}

func (g *valueGraphAsGraph[N, V]) PutEdge(nodeU, nodeV N) bool {
	return g.valueGraph.PutEdge(nodeU, nodeV)
}

func (g *valueGraphAsGraph[N, V]) RemoveNode(node N) bool {
	return g.valueGraph.RemoveNode(node)
}

func (g *valueGraphAsGraph[N, V]) RemoveEdge(nodeU, nodeV N) bool {
	return g.valueGraph.RemoveEdge(nodeU, nodeV)
}

func (g *valueGraphAsGraph[N, V]) Successors(node N) set.Set[N] {
	return g.valueGraph.Successors(node)
}

func (g *valueGraphAsGraph[N, V]) Predecessors(node N) set.Set[N] {
	return g.valueGraph.Predecessors(node)
}

func (g *valueGraphAsGraph[N, V]) AdjacentNodes(node N) set.Set[N] {
	return g.valueGraph.AdjacentNodes(node)
}

func (g *valueGraphAsGraph[N, V]) IncidentEdges(node N) set.Set[EndpointPair[N]] {
	return g.valueGraph.IncidentEdges(node)
}

func (g *valueGraphAsGraph[N, V]) Degree(node N) int {
	return g.valueGraph.Degree(node)
}

func (g *valueGraphAsGraph[N, V]) InDegree(node N) int {
	return g.valueGraph.InDegree(node)
}

func (g *valueGraphAsGraph[N, V]) OutDegree(node N) int {
	return g.valueGraph.OutDegree(node)
}

func (g *valueGraphAsGraph[N, V]) HasEdgeConnecting(nodeU, nodeV N) bool {
	return g.valueGraph.HasEdgeConnecting(nodeU, nodeV)
}

func (g *valueGraphAsGraph[N, V]) Size() int {
	return g.valueGraph.Size()
}

func (g *valueGraphAsGraph[N, V]) IsEmpty() bool {
	return g.valueGraph.IsEmpty()
}

func (g *valueGraphAsGraph[N, V]) Clear() {
	g.valueGraph.Clear()
}

func (g *valueGraphAsGraph[N, V]) Contains(node N) bool {
	return g.valueGraph.Contains(node)
}

func (g *valueGraphAsGraph[N, V]) ToSlice() []N {
	return g.valueGraph.ToSlice()
}

func (g *valueGraphAsGraph[N, V]) Iterator() common.Iterator[N] {
	return g.valueGraph.Iterator()
}

func (g *valueGraphAsGraph[N, V]) ForEach(fn func(N)) {
	g.valueGraph.ForEach(fn)
}

func (g *valueGraphAsGraph[N, V]) String() string {
	return g.valueGraph.String()
}