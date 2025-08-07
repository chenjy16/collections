package graph

// GraphBuilder provides a fluent interface for building graphs
type GraphBuilder[N comparable] struct {
	directed       bool
	allowSelfLoops bool
	nodeOrder      ElementOrder
}

// NewGraphBuilder creates a new graph builder
func NewGraphBuilder[N comparable]() *GraphBuilder[N] {
	return &GraphBuilder[N]{
		directed:       false,
		allowSelfLoops: false,
		nodeOrder:      Unordered,
	}
}

// Directed sets the graph to be directed
func (b *GraphBuilder[N]) Directed() *GraphBuilder[N] {
	newBuilder := *b
	newBuilder.directed = true
	return &newBuilder
}

// Undirected sets the graph to be undirected
func (b *GraphBuilder[N]) Undirected() *GraphBuilder[N] {
	newBuilder := *b
	newBuilder.directed = false
	return &newBuilder
}

// AllowSelfLoops allows self-loops in the graph
func (b *GraphBuilder[N]) AllowSelfLoops() *GraphBuilder[N] {
	newBuilder := *b
	newBuilder.allowSelfLoops = true
	return &newBuilder
}

// NodeOrder sets the node ordering
func (b *GraphBuilder[N]) NodeOrder(order ElementOrder) *GraphBuilder[N] {
	newBuilder := *b
	newBuilder.nodeOrder = order
	return &newBuilder
}

// Build creates a new mutable graph with the specified configuration
func (b *GraphBuilder[N]) Build() Graph[N] {
	return NewMutableGraph[N](b.directed, b.allowSelfLoops, b.nodeOrder)
}

// ValueGraphBuilder provides a fluent interface for building value graphs
type ValueGraphBuilder[N comparable, V any] struct {
	directed       bool
	allowSelfLoops bool
	nodeOrder      ElementOrder
}

// NewValueGraphBuilder creates a new value graph builder
func NewValueGraphBuilder[N comparable, V any]() *ValueGraphBuilder[N, V] {
	return &ValueGraphBuilder[N, V]{
		directed:       false,
		allowSelfLoops: false,
		nodeOrder:      Unordered,
	}
}

// Directed sets the graph to be directed
func (b *ValueGraphBuilder[N, V]) Directed() *ValueGraphBuilder[N, V] {
	newBuilder := *b
	newBuilder.directed = true
	return &newBuilder
}

// Undirected sets the graph to be undirected
func (b *ValueGraphBuilder[N, V]) Undirected() *ValueGraphBuilder[N, V] {
	newBuilder := *b
	newBuilder.directed = false
	return &newBuilder
}

// AllowSelfLoops allows self-loops in the graph
func (b *ValueGraphBuilder[N, V]) AllowSelfLoops() *ValueGraphBuilder[N, V] {
	newBuilder := *b
	newBuilder.allowSelfLoops = true
	return &newBuilder
}

// NodeOrder sets the node ordering
func (b *ValueGraphBuilder[N, V]) NodeOrder(order ElementOrder) *ValueGraphBuilder[N, V] {
	newBuilder := *b
	newBuilder.nodeOrder = order
	return &newBuilder
}

// Build creates a new mutable value graph with the specified configuration
func (b *ValueGraphBuilder[N, V]) Build() ValueGraph[N, V] {
	return NewMutableValueGraph[N, V](b.directed, b.allowSelfLoops, b.nodeOrder)
}

// NetworkBuilder provides a fluent interface for building networks
type NetworkBuilder[N comparable, E comparable] struct {
	directed           bool
	allowSelfLoops     bool
	allowParallelEdges bool
	nodeOrder          ElementOrder
	edgeOrder          ElementOrder
}

// NewNetworkBuilder creates a new network builder
func NewNetworkBuilder[N comparable, E comparable]() *NetworkBuilder[N, E] {
	return &NetworkBuilder[N, E]{
		directed:           false,
		allowSelfLoops:     false,
		allowParallelEdges: false,
		nodeOrder:          Unordered,
		edgeOrder:          Unordered,
	}
}

// Directed sets the network to be directed
func (b *NetworkBuilder[N, E]) Directed() *NetworkBuilder[N, E] {
	newBuilder := *b
	newBuilder.directed = true
	return &newBuilder
}

// Undirected sets the network to be undirected
func (b *NetworkBuilder[N, E]) Undirected() *NetworkBuilder[N, E] {
	newBuilder := *b
	newBuilder.directed = false
	return &newBuilder
}

// AllowSelfLoops allows self-loops in the network
func (b *NetworkBuilder[N, E]) AllowSelfLoops() *NetworkBuilder[N, E] {
	newBuilder := *b
	newBuilder.allowSelfLoops = true
	return &newBuilder
}

// AllowParallelEdges allows parallel edges in the network
func (b *NetworkBuilder[N, E]) AllowParallelEdges() *NetworkBuilder[N, E] {
	newBuilder := *b
	newBuilder.allowParallelEdges = true
	return &newBuilder
}

// NodeOrder sets the node ordering
func (b *NetworkBuilder[N, E]) NodeOrder(order ElementOrder) *NetworkBuilder[N, E] {
	newBuilder := *b
	newBuilder.nodeOrder = order
	return &newBuilder
}

// EdgeOrder sets the edge ordering
func (b *NetworkBuilder[N, E]) EdgeOrder(order ElementOrder) *NetworkBuilder[N, E] {
	newBuilder := *b
	newBuilder.edgeOrder = order
	return &newBuilder
}

// Build creates a new mutable network with the specified configuration
func (b *NetworkBuilder[N, E]) Build() Network[N, E] {
	return NewMutableNetwork[N, E](b.directed, b.allowSelfLoops, b.allowParallelEdges, b.nodeOrder, b.edgeOrder)
}

// Convenience factory functions

// DirectedGraph creates a new directed graph
func DirectedGraph[N comparable]() Graph[N] {
	return NewGraphBuilder[N]().Directed().Build()
}

// UndirectedGraph creates a new undirected graph
func UndirectedGraph[N comparable]() Graph[N] {
	return NewGraphBuilder[N]().Undirected().Build()
}

// DirectedValueGraph creates a new directed value graph
func DirectedValueGraph[N comparable, V any]() ValueGraph[N, V] {
	return NewValueGraphBuilder[N, V]().Directed().Build()
}

// UndirectedValueGraph creates a new undirected value graph
func UndirectedValueGraph[N comparable, V any]() ValueGraph[N, V] {
	return NewValueGraphBuilder[N, V]().Undirected().Build()
}

// DirectedNetwork creates a new directed network
func DirectedNetwork[N comparable, E comparable]() Network[N, E] {
	return NewNetworkBuilder[N, E]().Directed().Build()
}

// UndirectedNetwork creates a new undirected network
func UndirectedNetwork[N comparable, E comparable]() Network[N, E] {
	return NewNetworkBuilder[N, E]().Undirected().Build()
}

// DirectedMultigraph creates a new directed network that allows parallel edges
func DirectedMultigraph[N comparable, E comparable]() Network[N, E] {
	return NewNetworkBuilder[N, E]().Directed().AllowParallelEdges().Build()
}

// UndirectedMultigraph creates a new undirected network that allows parallel edges
func UndirectedMultigraph[N comparable, E comparable]() Network[N, E] {
	return NewNetworkBuilder[N, E]().Undirected().AllowParallelEdges().Build()
}