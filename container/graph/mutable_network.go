package graph

import (
	"fmt"
	"strings"

	"github.com/chenjianyu/collections/container/common"
	"github.com/chenjianyu/collections/container/set"
)

// MutableNetwork is a mutable implementation of Network
type MutableNetwork[N comparable, E comparable] struct {
	directed           bool
	allowSelfLoops     bool
	allowParallelEdges bool
	nodeOrder          ElementOrder
	edgeOrder          ElementOrder
	nodes              set.Set[N]
	edges              set.Set[E]
	// edgeToNodes maps each edge to its endpoint pair
	edgeToNodes map[E]EndpointPair[N]
	// nodeToEdges maps each node to its incident edges
	nodeToEdges map[N]set.Set[E]
	// For directed networks, we also maintain incoming and outgoing edges
	inEdges  map[N]set.Set[E]
	outEdges map[N]set.Set[E]
}

// NewMutableNetwork creates a new mutable network
func NewMutableNetwork[N comparable, E comparable](directed, allowSelfLoops, allowParallelEdges bool, nodeOrder, edgeOrder ElementOrder) *MutableNetwork[N, E] {
	return &MutableNetwork[N, E]{
		directed:           directed,
		allowSelfLoops:     allowSelfLoops,
		allowParallelEdges: allowParallelEdges,
		nodeOrder:          nodeOrder,
		edgeOrder:          edgeOrder,
		nodes:              set.New[N](),
		edges:              set.New[E](),
		edgeToNodes:        make(map[E]EndpointPair[N]),
		nodeToEdges:        make(map[N]set.Set[E]),
		inEdges:            make(map[N]set.Set[E]),
		outEdges:           make(map[N]set.Set[E]),
	}
}

// Nodes returns all nodes in this network
func (n *MutableNetwork[N, E]) Nodes() set.Set[N] {
	result := set.New[N]()
	n.nodes.ForEach(func(node N) {
		result.Add(node)
	})
	return result
}

// Edges returns all edges in this network
func (n *MutableNetwork[N, E]) Edges() set.Set[E] {
	result := set.New[E]()
	n.edges.ForEach(func(edge E) {
		result.Add(edge)
	})
	return result
}

// IsDirected returns true if this is a directed network
func (n *MutableNetwork[N, E]) IsDirected() bool {
	return n.directed
}

// AllowsSelfLoops returns true if this network allows self-loops
func (n *MutableNetwork[N, E]) AllowsSelfLoops() bool {
	return n.allowSelfLoops
}

// AllowsParallelEdges returns true if this network allows parallel edges
func (n *MutableNetwork[N, E]) AllowsParallelEdges() bool {
	return n.allowParallelEdges
}

// NodeOrder returns the ordering of nodes in this network
func (n *MutableNetwork[N, E]) NodeOrder() ElementOrder {
	return n.nodeOrder
}

// EdgeOrder returns the ordering of edges in this network
func (n *MutableNetwork[N, E]) EdgeOrder() ElementOrder {
	return n.edgeOrder
}

// AddNode adds a node to this network
func (n *MutableNetwork[N, E]) AddNode(node N) bool {
	if n.nodes.Contains(node) {
		return false
	}
	
	n.nodes.Add(node)
	n.nodeToEdges[node] = set.New[E]()
	if n.directed {
		n.inEdges[node] = set.New[E]()
		n.outEdges[node] = set.New[E]()
	}
	return true
}

// AddEdge adds an edge between two nodes
func (n *MutableNetwork[N, E]) AddEdge(edge E, nodeU, nodeV N) bool {
	// Check if edge already exists
	if n.edges.Contains(edge) {
		return false
	}
	
	// Check self-loop constraint
	if !n.allowSelfLoops && nodeU == nodeV {
		panic("Self-loops are not allowed in this network")
	}
	
	// Check parallel edge constraint
	if !n.allowParallelEdges && n.HasEdgeConnecting(nodeU, nodeV) {
		panic("Parallel edges are not allowed in this network")
	}
	
	// Add nodes if they don't exist
	n.AddNode(nodeU)
	n.AddNode(nodeV)
	
	// Add edge
	n.edges.Add(edge)
	n.edgeToNodes[edge] = NewEndpointPair(nodeU, nodeV)
	n.nodeToEdges[nodeU].Add(edge)
	n.nodeToEdges[nodeV].Add(edge)
	
	if n.directed {
		n.outEdges[nodeU].Add(edge)
		n.inEdges[nodeV].Add(edge)
	}
	
	return true
}

// RemoveNode removes a node and all its incident edges
func (n *MutableNetwork[N, E]) RemoveNode(node N) bool {
	if !n.nodes.Contains(node) {
		return false
	}
	
	// Remove all incident edges
	incidentEdges := n.nodeToEdges[node]
	edgesToRemove := incidentEdges.ToSlice()
	for _, edge := range edgesToRemove {
		n.RemoveEdge(edge)
	}
	
	// Remove the node itself
	n.nodes.Remove(node)
	delete(n.nodeToEdges, node)
	if n.directed {
		delete(n.inEdges, node)
		delete(n.outEdges, node)
	}
	
	return true
}

// RemoveEdge removes an edge
func (n *MutableNetwork[N, E]) RemoveEdge(edge E) bool {
	if !n.edges.Contains(edge) {
		return false
	}
	
	endpoints := n.edgeToNodes[edge]
	nodeU, nodeV := endpoints.NodeU, endpoints.NodeV
	
	// Remove edge from data structures
	n.edges.Remove(edge)
	delete(n.edgeToNodes, edge)
	n.nodeToEdges[nodeU].Remove(edge)
	n.nodeToEdges[nodeV].Remove(edge)
	
	if n.directed {
		n.outEdges[nodeU].Remove(edge)
		n.inEdges[nodeV].Remove(edge)
	}
	
	return true
}

// Successors returns the successors of a node
func (n *MutableNetwork[N, E]) Successors(node N) set.Set[N] {
	if !n.nodes.Contains(node) {
		panic(fmt.Sprintf("Node %v is not in the network", node))
	}
	
	result := set.New[N]()
	
	if n.directed {
		n.outEdges[node].ForEach(func(edge E) {
			endpoints := n.edgeToNodes[edge]
			result.Add(endpoints.NodeV)
		})
	} else {
		n.nodeToEdges[node].ForEach(func(edge E) {
			endpoints := n.edgeToNodes[edge]
			if endpoints.NodeU == node {
				result.Add(endpoints.NodeV)
			} else {
				result.Add(endpoints.NodeU)
			}
		})
	}
	
	return result
}

// Predecessors returns the predecessors of a node
func (n *MutableNetwork[N, E]) Predecessors(node N) set.Set[N] {
	if !n.nodes.Contains(node) {
		panic(fmt.Sprintf("Node %v is not in the network", node))
	}
	
	result := set.New[N]()
	
	if n.directed {
		n.inEdges[node].ForEach(func(edge E) {
			endpoints := n.edgeToNodes[edge]
			result.Add(endpoints.NodeU)
		})
	} else {
		// For undirected networks, predecessors are the same as successors
		return n.Successors(node)
	}
	
	return result
}

// AdjacentNodes returns all nodes adjacent to the given node
func (n *MutableNetwork[N, E]) AdjacentNodes(node N) set.Set[N] {
	if !n.nodes.Contains(node) {
		panic(fmt.Sprintf("Node %v is not in the network", node))
	}
	
	result := set.New[N]()
	
	n.nodeToEdges[node].ForEach(func(edge E) {
		endpoints := n.edgeToNodes[edge]
		if endpoints.NodeU == node {
			result.Add(endpoints.NodeV)
		} else {
			result.Add(endpoints.NodeU)
		}
	})
	
	return result
}

// IncidentEdges returns all edges incident to the given node
func (n *MutableNetwork[N, E]) IncidentEdges(node N) set.Set[E] {
	if !n.nodes.Contains(node) {
		panic(fmt.Sprintf("Node %v is not in the network", node))
	}
	
	result := set.New[E]()
	n.nodeToEdges[node].ForEach(func(edge E) {
		result.Add(edge)
	})
	return result
}

// InEdges returns all incoming edges to a node
func (n *MutableNetwork[N, E]) InEdges(node N) set.Set[E] {
	if !n.nodes.Contains(node) {
		panic(fmt.Sprintf("Node %v is not in the network", node))
	}
	
	result := set.New[E]()
	if n.directed {
		n.inEdges[node].ForEach(func(edge E) {
			result.Add(edge)
		})
	} else {
		// For undirected networks, all incident edges are both in and out
		return n.IncidentEdges(node)
	}
	return result
}

// OutEdges returns all outgoing edges from a node
func (n *MutableNetwork[N, E]) OutEdges(node N) set.Set[E] {
	if !n.nodes.Contains(node) {
		panic(fmt.Sprintf("Node %v is not in the network", node))
	}
	
	result := set.New[E]()
	if n.directed {
		n.outEdges[node].ForEach(func(edge E) {
			result.Add(edge)
		})
	} else {
		// For undirected networks, all incident edges are both in and out
		return n.IncidentEdges(node)
	}
	return result
}

// Degree returns the degree of a node
func (n *MutableNetwork[N, E]) Degree(node N) int {
	if !n.nodes.Contains(node) {
		panic(fmt.Sprintf("Node %v is not in the network", node))
	}
	
	return n.nodeToEdges[node].Size()
}

// InDegree returns the in-degree of a node
func (n *MutableNetwork[N, E]) InDegree(node N) int {
	if !n.nodes.Contains(node) {
		panic(fmt.Sprintf("Node %v is not in the network", node))
	}
	
	if n.directed {
		return n.inEdges[node].Size()
	}
	return n.nodeToEdges[node].Size()
}

// OutDegree returns the out-degree of a node
func (n *MutableNetwork[N, E]) OutDegree(node N) int {
	if !n.nodes.Contains(node) {
		panic(fmt.Sprintf("Node %v is not in the network", node))
	}
	
	if n.directed {
		return n.outEdges[node].Size()
	}
	return n.nodeToEdges[node].Size()
}

// IncidentNodes returns the nodes incident to an edge
func (n *MutableNetwork[N, E]) IncidentNodes(edge E) EndpointPair[N] {
	if !n.edges.Contains(edge) {
		panic(fmt.Sprintf("Edge %v is not in the network", edge))
	}
	
	return n.edgeToNodes[edge]
}

// AdjacentEdges returns all edges adjacent to the given edge
func (n *MutableNetwork[N, E]) AdjacentEdges(edge E) set.Set[E] {
	if !n.edges.Contains(edge) {
		panic(fmt.Sprintf("Edge %v is not in the network", edge))
	}
	
	result := set.New[E]()
	endpoints := n.edgeToNodes[edge]
	
	// Add all edges incident to the endpoints of this edge
	n.nodeToEdges[endpoints.NodeU].ForEach(func(incidentEdge E) {
		if incidentEdge != edge {
			result.Add(incidentEdge)
		}
	})
	
	n.nodeToEdges[endpoints.NodeV].ForEach(func(incidentEdge E) {
		if incidentEdge != edge {
			result.Add(incidentEdge)
		}
	})
	
	return result
}

// EdgesConnecting returns all edges connecting two nodes
func (n *MutableNetwork[N, E]) EdgesConnecting(nodeU, nodeV N) set.Set[E] {
	if !n.nodes.Contains(nodeU) || !n.nodes.Contains(nodeV) {
		return set.New[E]()
	}
	
	result := set.New[E]()
	
	n.nodeToEdges[nodeU].ForEach(func(edge E) {
		endpoints := n.edgeToNodes[edge]
		if (endpoints.NodeU == nodeU && endpoints.NodeV == nodeV) ||
			(!n.directed && endpoints.NodeU == nodeV && endpoints.NodeV == nodeU) {
			result.Add(edge)
		}
	})
	
	return result
}

// HasEdgeConnecting returns true if there's an edge between two nodes
func (n *MutableNetwork[N, E]) HasEdgeConnecting(nodeU, nodeV N) bool {
	return !n.EdgesConnecting(nodeU, nodeV).IsEmpty()
}

// AsGraph returns a view of this network as a basic graph
func (n *MutableNetwork[N, E]) AsGraph() Graph[N] {
	return &networkAsGraph[N, E]{n}
}

// Size returns the number of nodes in the network
func (n *MutableNetwork[N, E]) Size() int {
	return n.nodes.Size()
}

// IsEmpty returns true if the network has no nodes
func (n *MutableNetwork[N, E]) IsEmpty() bool {
	return n.nodes.IsEmpty()
}

// Clear removes all nodes and edges from the network
func (n *MutableNetwork[N, E]) Clear() {
	n.nodes.Clear()
	n.edges.Clear()
	n.edgeToNodes = make(map[E]EndpointPair[N])
	n.nodeToEdges = make(map[N]set.Set[E])
	n.inEdges = make(map[N]set.Set[E])
	n.outEdges = make(map[N]set.Set[E])
}

// Contains returns true if the network contains the given node
func (n *MutableNetwork[N, E]) Contains(node N) bool {
	return n.nodes.Contains(node)
}

// ToSlice returns all nodes as a slice
func (n *MutableNetwork[N, E]) ToSlice() []N {
	return n.nodes.ToSlice()
}

// Iterator returns an iterator over the nodes
func (n *MutableNetwork[N, E]) Iterator() common.Iterator[N] {
	return n.nodes.Iterator()
}

// ForEach applies a function to each node
func (n *MutableNetwork[N, E]) ForEach(fn func(N)) {
	n.nodes.ForEach(fn)
}

// String returns a string representation of the network
func (n *MutableNetwork[N, E]) String() string {
	var sb strings.Builder
	sb.WriteString("Network{")
	
	if n.directed {
		sb.WriteString("directed, ")
	} else {
		sb.WriteString("undirected, ")
	}
	
	sb.WriteString(fmt.Sprintf("nodes=%d, edges=%d", n.Size(), n.edges.Size()))
	sb.WriteString("}")
	
	return sb.String()
}

// networkAsGraph is an adapter that presents a Network as a Graph
type networkAsGraph[N comparable, E comparable] struct {
	network *MutableNetwork[N, E]
}

func (g *networkAsGraph[N, E]) Nodes() set.Set[N] {
	return g.network.Nodes()
}

func (g *networkAsGraph[N, E]) Edges() set.Set[EndpointPair[N]] {
	result := set.New[EndpointPair[N]]()
	g.network.edges.ForEach(func(edge E) {
		endpoints := g.network.edgeToNodes[edge]
		result.Add(endpoints)
	})
	return result
}

func (g *networkAsGraph[N, E]) IsDirected() bool {
	return g.network.IsDirected()
}

func (g *networkAsGraph[N, E]) AllowsSelfLoops() bool {
	return g.network.AllowsSelfLoops()
}

func (g *networkAsGraph[N, E]) NodeOrder() ElementOrder {
	return g.network.NodeOrder()
}

func (g *networkAsGraph[N, E]) AddNode(node N) bool {
	return g.network.AddNode(node)
}

func (g *networkAsGraph[N, E]) PutEdge(nodeU, nodeV N) bool {
	// For the graph view, we need to create a synthetic edge
	// This is a limitation of the adapter pattern
	panic("Cannot add edges through Graph view of Network - use Network.AddEdge instead")
}

func (g *networkAsGraph[N, E]) RemoveNode(node N) bool {
	return g.network.RemoveNode(node)
}

func (g *networkAsGraph[N, E]) RemoveEdge(nodeU, nodeV N) bool {
	// Remove all edges connecting these nodes
	edges := g.network.EdgesConnecting(nodeU, nodeV)
	if edges.IsEmpty() {
		return false
	}
	
	edges.ForEach(func(edge E) {
		g.network.RemoveEdge(edge)
	})
	return true
}

func (g *networkAsGraph[N, E]) Successors(node N) set.Set[N] {
	return g.network.Successors(node)
}

func (g *networkAsGraph[N, E]) Predecessors(node N) set.Set[N] {
	return g.network.Predecessors(node)
}

func (g *networkAsGraph[N, E]) AdjacentNodes(node N) set.Set[N] {
	return g.network.AdjacentNodes(node)
}

func (g *networkAsGraph[N, E]) IncidentEdges(node N) set.Set[EndpointPair[N]] {
	result := set.New[EndpointPair[N]]()
	g.network.IncidentEdges(node).ForEach(func(edge E) {
		endpoints := g.network.edgeToNodes[edge]
		result.Add(endpoints)
	})
	return result
}

func (g *networkAsGraph[N, E]) Degree(node N) int {
	return g.network.Degree(node)
}

func (g *networkAsGraph[N, E]) InDegree(node N) int {
	return g.network.InDegree(node)
}

func (g *networkAsGraph[N, E]) OutDegree(node N) int {
	return g.network.OutDegree(node)
}

func (g *networkAsGraph[N, E]) HasEdgeConnecting(nodeU, nodeV N) bool {
	return g.network.HasEdgeConnecting(nodeU, nodeV)
}

func (g *networkAsGraph[N, E]) Size() int {
	return g.network.Size()
}

func (g *networkAsGraph[N, E]) IsEmpty() bool {
	return g.network.IsEmpty()
}

func (g *networkAsGraph[N, E]) Clear() {
	g.network.Clear()
}

func (g *networkAsGraph[N, E]) Contains(node N) bool {
	return g.network.Contains(node)
}

func (g *networkAsGraph[N, E]) ToSlice() []N {
	return g.network.ToSlice()
}

func (g *networkAsGraph[N, E]) Iterator() common.Iterator[N] {
	return g.network.Iterator()
}

func (g *networkAsGraph[N, E]) ForEach(fn func(N)) {
	g.network.ForEach(fn)
}

func (g *networkAsGraph[N, E]) String() string {
	return g.network.String()
}