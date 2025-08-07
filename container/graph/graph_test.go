package graph

import (
	"testing"
)

func TestMutableGraphBasicOperations(t *testing.T) {
	// Test undirected graph
	g := UndirectedGraph[string]()
	
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
	
	// Test adding edges
	if !g.PutEdge("A", "B") {
		t.Error("Expected PutEdge to return true for new edge")
	}
	if g.PutEdge("A", "B") {
		t.Error("Expected PutEdge to return false for existing edge")
	}
	
	g.PutEdge("B", "C")
	
	// Test edge operations
	if !g.HasEdgeConnecting("A", "B") {
		t.Error("Expected edge between A and B")
	}
	if !g.HasEdgeConnecting("B", "A") {
		t.Error("Expected edge between B and A (undirected)")
	}
	
	// Test successors and predecessors
	successors := g.Successors("A")
	if !successors.Contains("B") {
		t.Error("Expected B to be successor of A")
	}
	
	predecessors := g.Predecessors("A")
	if !predecessors.Contains("B") {
		t.Error("Expected B to be predecessor of A (undirected)")
	}
	
	// Test adjacent nodes
	adjacent := g.AdjacentNodes("B")
	if !adjacent.Contains("A") || !adjacent.Contains("C") {
		t.Error("Expected A and C to be adjacent to B")
	}
	
	// Test degree
	if g.Degree("B") != 2 {
		t.Errorf("Expected degree 2 for B, got %d", g.Degree("B"))
	}
	
	// Test removing edge
	if !g.RemoveEdge("A", "B") {
		t.Error("Expected RemoveEdge to return true")
	}
	if !g.HasEdgeConnecting("A", "B") {
		// This should be false now
	} else {
		t.Error("Expected no edge between A and B after removal")
	}
	
	// Test removing node
	if !g.RemoveNode("C") {
		t.Error("Expected RemoveNode to return true")
	}
	if g.Contains("C") {
		t.Error("Expected C to be removed from graph")
	}
}

func TestDirectedGraph(t *testing.T) {
	g := DirectedGraph[int]()
	
	g.AddNode(1)
	g.AddNode(2)
	g.AddNode(3)
	
	g.PutEdge(1, 2)
	g.PutEdge(2, 3)
	
	// Test directed edge behavior
	if !g.HasEdgeConnecting(1, 2) {
		t.Error("Expected edge from 1 to 2")
	}
	if g.HasEdgeConnecting(2, 1) {
		t.Error("Expected no edge from 2 to 1 (directed)")
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

func TestGraphWithSelfLoops(t *testing.T) {
	g := NewGraphBuilder[string]().AllowSelfLoops().Build()
	
	g.AddNode("A")
	
	// Test self-loop
	if !g.PutEdge("A", "A") {
		t.Error("Expected to be able to add self-loop")
	}
	
	if !g.HasEdgeConnecting("A", "A") {
		t.Error("Expected self-loop to exist")
	}
}

func TestGraphPanicOnSelfLoop(t *testing.T) {
	g := UndirectedGraph[string]()
	g.AddNode("A")
	
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when adding self-loop to graph that doesn't allow them")
		}
	}()
	
	g.PutEdge("A", "A")
}

func TestGraphEdges(t *testing.T) {
	g := UndirectedGraph[string]()
	
	g.AddNode("A")
	g.AddNode("B")
	g.AddNode("C")
	
	g.PutEdge("A", "B")
	g.PutEdge("B", "C")
	
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

func TestGraphClear(t *testing.T) {
	g := UndirectedGraph[string]()
	
	g.AddNode("A")
	g.AddNode("B")
	g.PutEdge("A", "B")
	
	g.Clear()
	
	if !g.IsEmpty() {
		t.Error("Expected graph to be empty after clear")
	}
	
	if g.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", g.Size())
	}
}

func TestGraphIterator(t *testing.T) {
	g := UndirectedGraph[string]()
	
	nodes := []string{"A", "B", "C"}
	for _, node := range nodes {
		g.AddNode(node)
	}
	
	// Test Nodes method
	nodes_set := g.Nodes()
	if nodes_set.Size() != 3 {
		t.Errorf("Expected 3 nodes, got %d", nodes_set.Size())
	}
	
	for _, node := range nodes {
		if !nodes_set.Contains(node) {
			t.Errorf("Expected nodes set to contain %s", node)
		}
	}
}

func TestGraphString(t *testing.T) {
	g := DirectedGraph[string]()
	g.AddNode("A")
	g.AddNode("B")
	g.PutEdge("A", "B")
	
	str := g.String()
	if str == "" {
		t.Error("Expected non-empty string representation")
	}
	
	// Should contain information about being directed
	if len(str) < 5 {
		t.Error("Expected meaningful string representation")
	}
}