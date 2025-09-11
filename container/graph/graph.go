// Package graph provides graph data structure implementations
package graph

import (
	"github.com/chenjianyu/collections/container/common"
	"github.com/chenjianyu/collections/container/set"
)

// ElementOrder represents the ordering of elements in a graph
type ElementOrder int

const (
	// Unordered means elements have no ordering
	Unordered ElementOrder = iota
	// Insertion means elements are ordered by insertion order
	Insertion
	// Natural means elements are ordered by their natural ordering
	Natural
)

// EndpointPair represents a pair of endpoints for an edge
type EndpointPair[N comparable] struct {
	NodeU N
	NodeV N
}

// NewEndpointPair creates a new endpoint pair
func NewEndpointPair[N comparable](nodeU, nodeV N) EndpointPair[N] {
	return EndpointPair[N]{NodeU: nodeU, NodeV: nodeV}
}

// Graph represents a graph data structure with nodes and edges
// This is the basic graph interface similar to Guava's Graph
type Graph[N comparable] interface {
	// Nodes returns all nodes in this graph
	Nodes() set.Set[N]

	// Edges returns all edges in this graph as endpoint pairs
	Edges() set.Set[EndpointPair[N]]

	// IsDirected returns true if this is a directed graph
	IsDirected() bool

	// AllowsSelfLoops returns true if this graph allows self-loops
	AllowsSelfLoops() bool

	// NodeOrder returns the ordering of nodes in this graph
	NodeOrder() ElementOrder

	// AddNode adds a node to this graph
	// Returns true if the node was added, false if it already existed
	AddNode(node N) bool

	// PutEdge adds an edge between two nodes
	// Returns error if operation is not allowed (e.g., self-loops not allowed)
	PutEdge(nodeU, nodeV N) error

	// RemoveNode removes a node and all its incident edges
	// Returns true if the node was removed, false if it didn't exist
	RemoveNode(node N) bool

	// RemoveEdge removes an edge between two nodes
	// Returns true if the edge was removed, false if it didn't exist
	RemoveEdge(nodeU, nodeV N) bool

	// Successors returns the successors of a node
	// Returns error if node is not in the graph
	Successors(node N) (set.Set[N], error)

	// Predecessors returns the predecessors of a node
	// Returns error if node is not in the graph
	Predecessors(node N) (set.Set[N], error)

	// AdjacentNodes returns all nodes adjacent to the given node
	// Returns error if node is not in the graph
	AdjacentNodes(node N) (set.Set[N], error)

	// IncidentEdges returns all edges incident to the given node
	// Returns error if node is not in the graph
	IncidentEdges(node N) (set.Set[EndpointPair[N]], error)

	// Degree returns the degree of a node (number of incident edges)
	// Returns error if node is not in the graph
	Degree(node N) (int, error)

	// InDegree returns the in-degree of a node (for directed graphs)
	// Returns error if node is not in the graph
	InDegree(node N) (int, error)

	// OutDegree returns the out-degree of a node (for directed graphs)
	// Returns error if node is not in the graph
	OutDegree(node N) (int, error)

	// HasEdgeConnecting returns true if there's an edge between two nodes
	HasEdgeConnecting(nodeU, nodeV N) bool

	// Common container operations
	common.Container[N]
}

// ValueGraph represents a graph with values associated with edges
// Similar to Guava's ValueGraph
type ValueGraph[N comparable, V any] interface {
	// Graph operations
	Graph[N]

	// EdgeValue returns the value associated with an edge
	// Returns the value and true if the edge exists, zero value and false otherwise
	EdgeValue(nodeU, nodeV N) (V, bool)

	// EdgeValueOrDefault returns the value associated with an edge or a default value
	EdgeValueOrDefault(nodeU, nodeV N, defaultValue V) V

	// PutEdgeValue adds an edge with a value
	// Returns the previous value if the edge existed, zero value otherwise
	PutEdgeValue(nodeU, nodeV N, value V) (V, bool)

	// AsGraph returns a view of this value graph as a basic graph
	AsGraph() Graph[N]
}

// Network represents a graph with explicit edge objects
// Similar to Guava's Network
type Network[N comparable, E comparable] interface {
	// Nodes returns all nodes in this network
	Nodes() set.Set[N]

	// Edges returns all edges in this network
	Edges() set.Set[E]

	// IsDirected returns true if this is a directed network
	IsDirected() bool

	// AllowsSelfLoops returns true if this network allows self-loops
	AllowsSelfLoops() bool

	// AllowsParallelEdges returns true if this network allows parallel edges
	AllowsParallelEdges() bool

	// NodeOrder returns the ordering of nodes in this network
	NodeOrder() ElementOrder

	// EdgeOrder returns the ordering of edges in this network
	EdgeOrder() ElementOrder

	// AddNode adds a node to this network
	// Returns true if the node was added, false if it already existed
	AddNode(node N) bool

	// AddEdge adds an edge between two nodes
	// Returns error if operation is not allowed (e.g., self-loops or parallel edges not allowed)
	AddEdge(edge E, nodeU, nodeV N) error

	// RemoveNode removes a node and all its incident edges
	// Returns true if the node was removed, false if it didn't exist
	RemoveNode(node N) bool

	// RemoveEdge removes an edge
	// Returns true if the edge was removed, false if it didn't exist
	RemoveEdge(edge E) bool

	// Successors returns the successors of a node
	// Returns error if node is not in the network
	Successors(node N) (set.Set[N], error)

	// Predecessors returns the predecessors of a node
	// Returns error if node is not in the network
	Predecessors(node N) (set.Set[N], error)

	// AdjacentNodes returns all nodes adjacent to the given node
	// Returns error if node is not in the network
	AdjacentNodes(node N) (set.Set[N], error)

	// IncidentEdges returns all edges incident to the given node
	// Returns error if node is not in the network
	IncidentEdges(node N) (set.Set[E], error)

	// InEdges returns all incoming edges to a node (for directed networks)
	// Returns error if node is not in the network
	InEdges(node N) (set.Set[E], error)

	// OutEdges returns all outgoing edges from a node (for directed networks)
	// Returns error if node is not in the network
	OutEdges(node N) (set.Set[E], error)

	// Degree returns the degree of a node (number of incident edges)
	// Returns error if node is not in the network
	Degree(node N) (int, error)

	// InDegree returns the in-degree of a node (for directed networks)
	// Returns error if node is not in the network
	InDegree(node N) (int, error)

	// OutDegree returns the out-degree of a node (for directed networks)
	// Returns error if node is not in the network
	OutDegree(node N) (int, error)

	// IncidentNodes returns the nodes incident to an edge
	// Returns error if edge is not in the network
	IncidentNodes(edge E) (EndpointPair[N], error)

	// AdjacentEdges returns all edges adjacent to the given edge
	// Returns error if edge is not in the network
	AdjacentEdges(edge E) (set.Set[E], error)

	// EdgesConnecting returns all edges connecting two nodes
	EdgesConnecting(nodeU, nodeV N) set.Set[E]

	// HasEdgeConnecting returns true if there's an edge between two nodes
	HasEdgeConnecting(nodeU, nodeV N) bool

	// AsGraph returns a view of this network as a basic graph
	AsGraph() Graph[N]

	// Common container operations
	common.Container[N]
}
