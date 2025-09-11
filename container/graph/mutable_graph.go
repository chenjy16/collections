package graph

import (
	"fmt"
	"strings"

	"github.com/chenjianyu/collections/container/common"
	"github.com/chenjianyu/collections/container/set"
)

// MutableGraph is a mutable implementation of Graph
type MutableGraph[N comparable] struct {
	directed       bool
	allowSelfLoops bool
	nodeOrder      ElementOrder
	nodes          set.Set[N]
	// adjacencyMap maps each node to its adjacent nodes
	adjacencyMap map[N]set.Set[N]
	// For directed graphs, we also maintain predecessors
	predecessorMap map[N]set.Set[N]
}

// NewMutableGraph creates a new mutable graph
func NewMutableGraph[N comparable](directed, allowSelfLoops bool, nodeOrder ElementOrder) *MutableGraph[N] {
	return &MutableGraph[N]{
		directed:       directed,
		allowSelfLoops: allowSelfLoops,
		nodeOrder:      nodeOrder,
		nodes:          set.New[N](),
		adjacencyMap:   make(map[N]set.Set[N]),
		predecessorMap: make(map[N]set.Set[N]),
	}
}

// Nodes returns all nodes in this graph
func (g *MutableGraph[N]) Nodes() set.Set[N] {
	result := set.New[N]()
	g.nodes.ForEach(func(node N) {
		result.Add(node)
	})
	return result
}

// Edges returns all edges in this graph as endpoint pairs
func (g *MutableGraph[N]) Edges() set.Set[EndpointPair[N]] {
	result := set.New[EndpointPair[N]]()

	for node, successors := range g.adjacencyMap {
		successors.ForEach(func(successor N) {
			if g.directed {
				result.Add(NewEndpointPair(node, successor))
			} else {
				// For undirected graphs, only add each edge once
				// Use a consistent ordering to avoid duplicates
				if fmt.Sprintf("%v", node) <= fmt.Sprintf("%v", successor) {
					result.Add(NewEndpointPair(node, successor))
				}
			}
		})
	}

	return result
}

// IsDirected returns true if this is a directed graph
func (g *MutableGraph[N]) IsDirected() bool {
	return g.directed
}

// AllowsSelfLoops returns true if this graph allows self-loops
func (g *MutableGraph[N]) AllowsSelfLoops() bool {
	return g.allowSelfLoops
}

// NodeOrder returns the ordering of nodes in this graph
func (g *MutableGraph[N]) NodeOrder() ElementOrder {
	return g.nodeOrder
}

// AddNode adds a node to this graph
func (g *MutableGraph[N]) AddNode(node N) bool {
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
func (g *MutableGraph[N]) PutEdge(nodeU, nodeV N) error {
	// Check self-loop constraint
	if !g.allowSelfLoops && nodeU == nodeV {
		return common.SelfLoopNotAllowedError(nodeU)
	}

	// Add nodes if they don't exist
	g.AddNode(nodeU)
	g.AddNode(nodeV)

	// Check if edge already exists
	if g.adjacencyMap[nodeU].Contains(nodeV) {
		return nil // Edge already exists, no error
	}

	// Add edge
	g.adjacencyMap[nodeU].Add(nodeV)
	if g.directed {
		g.predecessorMap[nodeV].Add(nodeU)
	} else {
		// For undirected graphs, add both directions
		g.adjacencyMap[nodeV].Add(nodeU)
	}

	return nil
}

// RemoveNode removes a node and all its incident edges
func (g *MutableGraph[N]) RemoveNode(node N) bool {
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
func (g *MutableGraph[N]) RemoveEdge(nodeU, nodeV N) bool {
	if !g.nodes.Contains(nodeU) || !g.nodes.Contains(nodeV) {
		return false
	}

	if !g.adjacencyMap[nodeU].Contains(nodeV) {
		return false
	}

	g.adjacencyMap[nodeU].Remove(nodeV)
	if g.directed {
		g.predecessorMap[nodeV].Remove(nodeU)
	} else {
		g.adjacencyMap[nodeV].Remove(nodeU)
	}

	return true
}

// Successors returns the successors of a node
func (g *MutableGraph[N]) Successors(node N) (set.Set[N], error) {
	if !g.nodes.Contains(node) {
		return nil, common.NodeNotFoundError(node)
	}

	result := set.New[N]()
	g.adjacencyMap[node].ForEach(func(successor N) {
		result.Add(successor)
	})
	return result, nil
}

// Predecessors returns the predecessors of a node
func (g *MutableGraph[N]) Predecessors(node N) (set.Set[N], error) {
	if !g.nodes.Contains(node) {
		return nil, common.NodeNotFoundError(node)
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
	return result, nil
}

// AdjacentNodes returns all nodes adjacent to the given node
func (g *MutableGraph[N]) AdjacentNodes(node N) (set.Set[N], error) {
	if !g.nodes.Contains(node) {
		return nil, common.NodeNotFoundError(node)
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

	return result, nil
}

// IncidentEdges returns all edges incident to the given node
func (g *MutableGraph[N]) IncidentEdges(node N) (set.Set[EndpointPair[N]], error) {
	if !g.nodes.Contains(node) {
		return nil, common.NodeNotFoundError(node)
	}

	result := set.New[EndpointPair[N]]()

	// Add outgoing edges
	g.adjacencyMap[node].ForEach(func(successor N) {
		result.Add(NewEndpointPair(node, successor))
	})

	// Add incoming edges (for directed graphs)
	if g.directed {
		g.predecessorMap[node].ForEach(func(predecessor N) {
			result.Add(NewEndpointPair(predecessor, node))
		})
	}

	return result, nil
}

// Degree returns the degree of a node
func (g *MutableGraph[N]) Degree(node N) (int, error) {
	if !g.nodes.Contains(node) {
		return 0, common.NodeNotFoundError(node)
	}

	if g.directed {
		inDegree, _ := g.InDegree(node)
		outDegree, _ := g.OutDegree(node)
		return inDegree + outDegree, nil
	}
	return g.adjacencyMap[node].Size(), nil
}

// InDegree returns the in-degree of a node
func (g *MutableGraph[N]) InDegree(node N) (int, error) {
	if !g.nodes.Contains(node) {
		return 0, common.NodeNotFoundError(node)
	}

	if g.directed {
		return g.predecessorMap[node].Size(), nil
	}
	return g.adjacencyMap[node].Size(), nil
}

// OutDegree returns the out-degree of a node
func (g *MutableGraph[N]) OutDegree(node N) (int, error) {
	if !g.nodes.Contains(node) {
		return 0, common.NodeNotFoundError(node)
	}

	return g.adjacencyMap[node].Size(), nil
}

// HasEdgeConnecting returns true if there's an edge between two nodes
func (g *MutableGraph[N]) HasEdgeConnecting(nodeU, nodeV N) bool {
	if !g.nodes.Contains(nodeU) || !g.nodes.Contains(nodeV) {
		return false
	}

	return g.adjacencyMap[nodeU].Contains(nodeV)
}

// Size returns the number of nodes in the graph
func (g *MutableGraph[N]) Size() int {
	return g.nodes.Size()
}

// IsEmpty returns true if the graph has no nodes
func (g *MutableGraph[N]) IsEmpty() bool {
	return g.nodes.IsEmpty()
}

// Clear removes all nodes and edges from the graph
func (g *MutableGraph[N]) Clear() {
	g.nodes.Clear()
	g.adjacencyMap = make(map[N]set.Set[N])
	g.predecessorMap = make(map[N]set.Set[N])
}

// Contains returns true if the graph contains the given node
func (g *MutableGraph[N]) Contains(node N) bool {
	return g.nodes.Contains(node)
}

// ToSlice returns all nodes as a slice
func (g *MutableGraph[N]) ToSlice() []N {
	return g.nodes.ToSlice()
}

// Iterator returns an iterator over the nodes
func (g *MutableGraph[N]) Iterator() common.Iterator[N] {
	return g.nodes.Iterator()
}

// ForEach applies a function to each node
func (g *MutableGraph[N]) ForEach(fn func(N)) {
	g.nodes.ForEach(fn)
}

// String returns a string representation of the graph
func (g *MutableGraph[N]) String() string {
	var sb strings.Builder
	sb.WriteString("Graph{")

	if g.directed {
		sb.WriteString("directed, ")
	} else {
		sb.WriteString("undirected, ")
	}

	sb.WriteString(fmt.Sprintf("nodes=%d, edges=%d", g.Size(), g.Edges().Size()))
	sb.WriteString("}")

	return sb.String()
}
