package graph

import (
	"testing"
)

func TestGraphBuilder(t *testing.T) {
	// Test default undirected graph
	graph := NewGraphBuilder[string]().Build()
	
	if graph.IsDirected() {
		t.Error("Expected undirected graph by default")
	}
	
	if graph.AllowsSelfLoops() {
		t.Error("Expected self-loops to be disallowed by default")
	}
	
	// Test directed graph
	directedGraph := NewGraphBuilder[string]().Directed().Build()
	
	if !directedGraph.IsDirected() {
		t.Error("Expected directed graph")
	}
	
	// Test self-loops allowed
	selfLoopGraph := NewGraphBuilder[string]().AllowSelfLoops().Build()
	
	if !selfLoopGraph.AllowsSelfLoops() {
		t.Error("Expected self-loops to be allowed")
	}
	
	// Test chaining
	complexGraph := NewGraphBuilder[int]().
		Directed().
		AllowSelfLoops().
		Build()
	
	if !complexGraph.IsDirected() {
		t.Error("Expected directed graph")
	}
	
	if !complexGraph.AllowsSelfLoops() {
		t.Error("Expected self-loops to be allowed")
	}
}

func TestValueGraphBuilder(t *testing.T) {
	// Test default undirected value graph
	valueGraph := NewValueGraphBuilder[string, int]().Build()
	
	if valueGraph.IsDirected() {
		t.Error("Expected undirected value graph by default")
	}
	
	if valueGraph.AllowsSelfLoops() {
		t.Error("Expected self-loops to be disallowed by default")
	}
	
	// Test directed value graph
	directedValueGraph := NewValueGraphBuilder[string, int]().Directed().Build()
	
	if !directedValueGraph.IsDirected() {
		t.Error("Expected directed value graph")
	}
	
	// Test self-loops allowed
	selfLoopValueGraph := NewValueGraphBuilder[string, int]().AllowSelfLoops().Build()
	
	if !selfLoopValueGraph.AllowsSelfLoops() {
		t.Error("Expected self-loops to be allowed")
	}
	
	// Test chaining
	complexValueGraph := NewValueGraphBuilder[int, string]().
		Directed().
		AllowSelfLoops().
		Build()
	
	if !complexValueGraph.IsDirected() {
		t.Error("Expected directed value graph")
	}
	
	if !complexValueGraph.AllowsSelfLoops() {
		t.Error("Expected self-loops to be allowed")
	}
}

func TestNetworkBuilder(t *testing.T) {
	// Test default undirected network
	network := NewNetworkBuilder[string, string]().Build()
	
	if network.IsDirected() {
		t.Error("Expected undirected network by default")
	}
	
	if network.AllowsSelfLoops() {
		t.Error("Expected self-loops to be disallowed by default")
	}
	
	if network.AllowsParallelEdges() {
		t.Error("Expected parallel edges to be disallowed by default")
	}
	
	// Test directed network
	directedNetwork := NewNetworkBuilder[string, string]().Directed().Build()
	
	if !directedNetwork.IsDirected() {
		t.Error("Expected directed network")
	}
	
	// Test self-loops allowed
	selfLoopNetwork := NewNetworkBuilder[string, string]().AllowSelfLoops().Build()
	
	if !selfLoopNetwork.AllowsSelfLoops() {
		t.Error("Expected self-loops to be allowed")
	}
	
	// Test parallel edges allowed
	parallelEdgeNetwork := NewNetworkBuilder[string, string]().AllowParallelEdges().Build()
	
	if !parallelEdgeNetwork.AllowsParallelEdges() {
		t.Error("Expected parallel edges to be allowed")
	}
	
	// Test chaining
	complexNetwork := NewNetworkBuilder[int, string]().
		Directed().
		AllowSelfLoops().
		AllowParallelEdges().
		Build()
	
	if !complexNetwork.IsDirected() {
		t.Error("Expected directed network")
	}
	
	if !complexNetwork.AllowsSelfLoops() {
		t.Error("Expected self-loops to be allowed")
	}
	
	if !complexNetwork.AllowsParallelEdges() {
		t.Error("Expected parallel edges to be allowed")
	}
}

func TestConvenienceFunctions(t *testing.T) {
	// Test UndirectedGraph
	undirectedGraph := UndirectedGraph[string]()
	if undirectedGraph.IsDirected() {
		t.Error("Expected undirected graph")
	}
	
	// Test DirectedGraph
	directedGraph := DirectedGraph[string]()
	if !directedGraph.IsDirected() {
		t.Error("Expected directed graph")
	}
	
	// Test UndirectedValueGraph
	undirectedValueGraph := UndirectedValueGraph[string, int]()
	if undirectedValueGraph.IsDirected() {
		t.Error("Expected undirected value graph")
	}
	
	// Test DirectedValueGraph
	directedValueGraph := DirectedValueGraph[string, int]()
	if !directedValueGraph.IsDirected() {
		t.Error("Expected directed value graph")
	}
	
	// Test UndirectedNetwork
	undirectedNetwork := UndirectedNetwork[string, string]()
	if undirectedNetwork.IsDirected() {
		t.Error("Expected undirected network")
	}
	
	// Test DirectedNetwork
	directedNetwork := DirectedNetwork[string, string]()
	if !directedNetwork.IsDirected() {
		t.Error("Expected directed network")
	}
}

func TestBuilderImmutability(t *testing.T) {
	// Test that builder methods return new instances
	builder1 := NewGraphBuilder[string]()
	builder2 := builder1.Directed()
	
	// Original builder should still create undirected graphs
	graph1 := builder1.Build()
	if graph1.IsDirected() {
		t.Error("Expected original builder to still create undirected graphs")
	}
	
	// New builder should create directed graphs
	graph2 := builder2.Build()
	if !graph2.IsDirected() {
		t.Error("Expected new builder to create directed graphs")
	}
}

func TestBuilderReuse(t *testing.T) {
	// Test that builders can be reused
	builder := NewGraphBuilder[int]().Directed().AllowSelfLoops()
	
	graph1 := builder.Build()
	graph2 := builder.Build()
	
	// Both graphs should have the same properties
	if !graph1.IsDirected() || !graph2.IsDirected() {
		t.Error("Expected both graphs to be directed")
	}
	
	if !graph1.AllowsSelfLoops() || !graph2.AllowsSelfLoops() {
		t.Error("Expected both graphs to allow self-loops")
	}
	
	// But they should be different instances
	graph1.AddNode(1)
	if graph2.Contains(1) {
		t.Error("Expected graphs to be independent instances")
	}
}

func TestValueGraphBuilderTypes(t *testing.T) {
	// Test with different value types
	stringValueGraph := NewValueGraphBuilder[string, string]().Build()
	intValueGraph := NewValueGraphBuilder[int, int]().Build()
	mixedValueGraph := NewValueGraphBuilder[string, int]().Build()
	
	// Test basic operations to ensure types work correctly
	stringValueGraph.AddNode("A")
	stringValueGraph.AddNode("B")
	stringValueGraph.PutEdgeValue("A", "B", "edge_value")
	
	intValueGraph.AddNode(1)
	intValueGraph.AddNode(2)
	intValueGraph.PutEdgeValue(1, 2, 42)
	
	mixedValueGraph.AddNode("node1")
	mixedValueGraph.AddNode("node2")
	mixedValueGraph.PutEdgeValue("node1", "node2", 100)
	
	// Verify values
	if val, exists := stringValueGraph.EdgeValue("A", "B"); !exists || val != "edge_value" {
		t.Error("Expected string edge value to be 'edge_value'")
	}
	
	if val, exists := intValueGraph.EdgeValue(1, 2); !exists || val != 42 {
		t.Error("Expected int edge value to be 42")
	}
	
	if val, exists := mixedValueGraph.EdgeValue("node1", "node2"); !exists || val != 100 {
		t.Error("Expected mixed edge value to be 100")
	}
}

func TestNetworkBuilderTypes(t *testing.T) {
	// Test with different edge types
	stringNetwork := NewNetworkBuilder[string, string]().Build()
	intNetwork := NewNetworkBuilder[int, int]().Build()
	mixedNetwork := NewNetworkBuilder[string, int]().Build()
	
	// Test basic operations to ensure types work correctly
	stringNetwork.AddNode("A")
	stringNetwork.AddNode("B")
	stringNetwork.AddEdge("edge1", "A", "B")
	
	intNetwork.AddNode(1)
	intNetwork.AddNode(2)
	intNetwork.AddEdge(100, 1, 2)
	
	mixedNetwork.AddNode("node1")
	mixedNetwork.AddNode("node2")
	mixedNetwork.AddEdge(42, "node1", "node2")
	
	// Verify edges exist
	if !stringNetwork.HasEdgeConnecting("A", "B") {
		t.Error("Expected string network to have edge A-B")
	}
	
	if !intNetwork.HasEdgeConnecting(1, 2) {
		t.Error("Expected int network to have edge 1-2")
	}
	
	if !mixedNetwork.HasEdgeConnecting("node1", "node2") {
		t.Error("Expected mixed network to have edge node1-node2")
	}
}