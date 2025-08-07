package graph

import (
	"testing"
)

func TestMutableValueGraphBasicOperations(t *testing.T) {
	// Test undirected value graph
	g := UndirectedValueGraph[string, int]()
	
	// Test adding nodes
	if !g.AddNode("A") {
		t.Error("Expected AddNode to return true for new node")
	}
	if g.AddNode("A") {
		t.Error("Expected AddNode to return false for existing node")
	}
	
	g.AddNode("B")
	g.AddNode("C")
	
	// Test node operations
	if g.Size() != 3 {
		t.Errorf("Expected size 3, got %d", g.Size())
	}
	
	if !g.Contains("A") {
		t.Error("Expected graph to contain node A")
	}
	
	// Test adding edges with values
	prevValue, existed := g.PutEdgeValue("A", "B", 10)
	if existed {
		t.Error("Expected PutEdgeValue to return false for new edge")
	}
	if prevValue != 0 {
		t.Errorf("Expected previous value 0, got %d", prevValue)
	}
	
	// Test updating edge value
	prevValue, existed = g.PutEdgeValue("A", "B", 20)
	if !existed {
		t.Error("Expected PutEdgeValue to return true for existing edge")
	}
	if prevValue != 10 {
		t.Errorf("Expected previous value 10, got %d", prevValue)
	}
	
	g.PutEdgeValue("B", "C", 30)
	
	// Test edge value operations
	value, exists := g.EdgeValue("A", "B")
	if !exists {
		t.Error("Expected edge value to exist")
	}
	if value != 20 {
		t.Errorf("Expected edge value 20, got %d", value)
	}
	
	// Test edge value for undirected graph (reverse direction)
	value, exists = g.EdgeValue("B", "A")
	if !exists {
		t.Error("Expected edge value to exist in reverse direction (undirected)")
	}
	if value != 20 {
		t.Errorf("Expected edge value 20 in reverse direction, got %d", value)
	}
	
	// Test EdgeValueOrDefault
	defaultValue := g.EdgeValueOrDefault("A", "C", 999)
	if defaultValue != 999 {
		t.Errorf("Expected default value 999, got %d", defaultValue)
	}
	
	existingValue := g.EdgeValueOrDefault("A", "B", 999)
	if existingValue != 20 {
		t.Errorf("Expected existing value 20, got %d", existingValue)
	}
	
	// Test basic graph operations
	if !g.HasEdgeConnecting("A", "B") {
		t.Error("Expected edge between A and B")
	}
	
	// Test successors and predecessors
	successors := g.Successors("A")
	if !successors.Contains("B") {
		t.Error("Expected B to be successor of A")
	}
	
	// Test AsGraph
	basicGraph := g.AsGraph()
	if !basicGraph.HasEdgeConnecting("A", "B") {
		t.Error("Expected basic graph view to have edge A-B")
	}
	
	// Test removing edge
	if !g.RemoveEdge("A", "B") {
		t.Error("Expected RemoveEdge to return true")
	}
	
	_, exists = g.EdgeValue("A", "B")
	if exists {
		t.Error("Expected edge value to not exist after removal")
	}
}

func TestDirectedValueGraph(t *testing.T) {
	g := DirectedValueGraph[int, string]()
	
	g.AddNode(1)
	g.AddNode(2)
	g.AddNode(3)
	
	g.PutEdgeValue(1, 2, "edge1-2")
	g.PutEdgeValue(2, 3, "edge2-3")
	
	// Test directed edge behavior
	value, exists := g.EdgeValue(1, 2)
	if !exists || value != "edge1-2" {
		t.Error("Expected edge value from 1 to 2")
	}
	
	_, exists = g.EdgeValue(2, 1)
	if exists {
		t.Error("Expected no edge value from 2 to 1 (directed)")
	}
	
	// Test successors and predecessors
	successors := g.Successors(1)
	if !successors.Contains(2) {
		t.Error("Expected 2 to be successor of 1")
	}
	
	predecessors := g.Predecessors(2)
	if !predecessors.Contains(1) {
		t.Error("Expected 1 to be predecessor of 2")
	}
	
	// Test in-degree and out-degree
	if g.InDegree(2) != 1 {
		t.Errorf("Expected in-degree 1 for node 2, got %d", g.InDegree(2))
	}
	if g.OutDegree(2) != 1 {
		t.Errorf("Expected out-degree 1 for node 2, got %d", g.OutDegree(2))
	}
}

func TestValueGraphWithSelfLoops(t *testing.T) {
	g := NewValueGraphBuilder[string, float64]().AllowSelfLoops().Build()
	
	g.AddNode("A")
	
	// Test self-loop with value
	g.PutEdgeValue("A", "A", 3.14)
	
	value, exists := g.EdgeValue("A", "A")
	if !exists || value != 3.14 {
		t.Error("Expected self-loop with value 3.14")
	}
}

func TestValueGraphPanicOnSelfLoop(t *testing.T) {
	g := UndirectedValueGraph[string, int]()
	g.AddNode("A")
	
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when adding self-loop to graph that doesn't allow them")
		}
	}()
	
	g.PutEdgeValue("A", "A", 42)
}

func TestValueGraphEdges(t *testing.T) {
	g := UndirectedValueGraph[string, int]()
	
	g.AddNode("A")
	g.AddNode("B")
	g.AddNode("C")
	
	g.PutEdgeValue("A", "B", 10)
	g.PutEdgeValue("B", "C", 20)
	
	edges := g.Edges()
	if edges.Size() != 2 {
		t.Errorf("Expected 2 edges, got %d", edges.Size())
	}
	
	// Test incident edges
	incidentEdges := g.IncidentEdges("B")
	if incidentEdges.Size() != 2 {
		t.Errorf("Expected 2 incident edges for B, got %d", incidentEdges.Size())
	}
}

func TestValueGraphClear(t *testing.T) {
	g := UndirectedValueGraph[string, int]()
	
	g.AddNode("A")
	g.AddNode("B")
	g.PutEdgeValue("A", "B", 42)
	
	g.Clear()
	
	if !g.IsEmpty() {
		t.Error("Expected graph to be empty after clear")
	}
	
	if g.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", g.Size())
	}
	
	_, exists := g.EdgeValue("A", "B")
	if exists {
		t.Error("Expected no edge values after clear")
	}
}

func TestValueGraphString(t *testing.T) {
	g := DirectedValueGraph[string, int]()
	g.AddNode("A")
	g.AddNode("B")
	g.PutEdgeValue("A", "B", 42)
	
	str := g.String()
	if str == "" {
		t.Error("Expected non-empty string representation")
	}
	
	// Should contain information about being directed
	if len(str) < 5 {
		t.Error("Expected meaningful string representation")
	}
}

func TestValueGraphAsGraphView(t *testing.T) {
	valueGraph := UndirectedValueGraph[string, int]()
	valueGraph.AddNode("A")
	valueGraph.AddNode("B")
	valueGraph.PutEdgeValue("A", "B", 42)
	
	// Get graph view
	graph := valueGraph.AsGraph()
	
	// Test that graph view reflects the value graph
	if !graph.HasEdgeConnecting("A", "B") {
		t.Error("Expected graph view to have edge A-B")
	}
	
	if graph.Size() != 2 {
		t.Errorf("Expected graph view to have 2 nodes, got %d", graph.Size())
	}
	
	// Test that modifications through graph view affect value graph
	graph.AddNode("C")
	if !valueGraph.Contains("C") {
		t.Error("Expected value graph to contain node C added through graph view")
	}
	
	// Test that basic graph operations work
	successors := graph.Successors("A")
	if !successors.Contains("B") {
		t.Error("Expected B to be successor of A in graph view")
	}
}