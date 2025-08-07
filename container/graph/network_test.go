package graph

import (
	"testing"
)

func TestMutableNetworkBasicOperations(t *testing.T) {
	// Test undirected network
	n := UndirectedNetwork[string, string]()
	
	// Test adding nodes
	if !n.AddNode("A") {
		t.Error("Expected AddNode to return true for new node")
	}
	if n.AddNode("A") {
		t.Error("Expected AddNode to return false for existing node")
	}
	
	n.AddNode("B")
	n.AddNode("C")
	
	// Test node operations
	if n.Size() != 3 {
		t.Errorf("Expected size 3, got %d", n.Size())
	}
	
	if !n.Contains("A") {
		t.Error("Expected network to contain node A")
	}
	
	// Test adding edges
	if !n.AddEdge("edge1", "A", "B") {
		t.Error("Expected AddEdge to return true for new edge")
	}
	if n.AddEdge("edge1", "A", "B") {
		t.Error("Expected AddEdge to return false for existing edge")
	}
	
	n.AddEdge("edge2", "B", "C")
	
	// Test edge operations
	if !n.HasEdgeConnecting("A", "B") {
		t.Error("Expected edge between A and B")
	}
	if !n.HasEdgeConnecting("B", "A") {
		t.Error("Expected edge between B and A (undirected)")
	}
	
	// Test edges connecting
	edges := n.EdgesConnecting("A", "B")
	if edges.Size() != 1 {
		t.Errorf("Expected 1 edge connecting A and B, got %d", edges.Size())
	}
	if !edges.Contains("edge1") {
		t.Error("Expected edge1 to connect A and B")
	}
	
	// Test incident nodes
	endpoints := n.IncidentNodes("edge1")
	if (endpoints.NodeU != "A" || endpoints.NodeV != "B") && 
	   (endpoints.NodeU != "B" || endpoints.NodeV != "A") {
		t.Error("Expected edge1 to connect A and B")
	}
	
	// Test successors and predecessors
	successors := n.Successors("A")
	if !successors.Contains("B") {
		t.Error("Expected B to be successor of A")
	}
	
	predecessors := n.Predecessors("A")
	if !predecessors.Contains("B") {
		t.Error("Expected B to be predecessor of A (undirected)")
	}
	
	// Test adjacent nodes
	adjacent := n.AdjacentNodes("B")
	if !adjacent.Contains("A") || !adjacent.Contains("C") {
		t.Error("Expected A and C to be adjacent to B")
	}
	
	// Test incident edges
	incidentEdges := n.IncidentEdges("B")
	if incidentEdges.Size() != 2 {
		t.Errorf("Expected 2 incident edges for B, got %d", incidentEdges.Size())
	}
	if !incidentEdges.Contains("edge1") || !incidentEdges.Contains("edge2") {
		t.Error("Expected edge1 and edge2 to be incident to B")
	}
	
	// Test degree
	if n.Degree("B") != 2 {
		t.Errorf("Expected degree 2 for B, got %d", n.Degree("B"))
	}
	
	// Test removing edge
	if !n.RemoveEdge("edge1") {
		t.Error("Expected RemoveEdge to return true")
	}
	if n.HasEdgeConnecting("A", "B") {
		t.Error("Expected no edge between A and B after removal")
	}
	
	// Test removing node
	if !n.RemoveNode("C") {
		t.Error("Expected RemoveNode to return true")
	}
	if n.Contains("C") {
		t.Error("Expected C to be removed from network")
	}
	if n.Edges().Contains("edge2") {
		t.Error("Expected edge2 to be removed when C was removed")
	}
}

func TestDirectedNetwork(t *testing.T) {
	n := DirectedNetwork[int, string]()
	
	n.AddNode(1)
	n.AddNode(2)
	n.AddNode(3)
	
	n.AddEdge("e1-2", 1, 2)
	n.AddEdge("e2-3", 2, 3)
	
	// Test directed edge behavior
	if !n.HasEdgeConnecting(1, 2) {
		t.Error("Expected edge from 1 to 2")
	}
	if n.HasEdgeConnecting(2, 1) {
		t.Error("Expected no edge from 2 to 1 (directed)")
	}
	
	// Test successors and predecessors
	successors := n.Successors(1)
	if !successors.Contains(2) {
		t.Error("Expected 2 to be successor of 1")
	}
	
	predecessors := n.Predecessors(2)
	if !predecessors.Contains(1) {
		t.Error("Expected 1 to be predecessor of 2")
	}
	
	// Test in-edges and out-edges
	inEdges := n.InEdges(2)
	if !inEdges.Contains("e1-2") {
		t.Error("Expected e1-2 to be in-edge of 2")
	}
	
	outEdges := n.OutEdges(2)
	if !outEdges.Contains("e2-3") {
		t.Error("Expected e2-3 to be out-edge of 2")
	}
	
	// Test in-degree and out-degree
	if n.InDegree(2) != 1 {
		t.Errorf("Expected in-degree 1 for node 2, got %d", n.InDegree(2))
	}
	if n.OutDegree(2) != 1 {
		t.Errorf("Expected out-degree 1 for node 2, got %d", n.OutDegree(2))
	}
}

func TestNetworkWithSelfLoops(t *testing.T) {
	n := NewNetworkBuilder[string, string]().AllowSelfLoops().Build()
	
	n.AddNode("A")
	
	// Test self-loop
	if !n.AddEdge("self", "A", "A") {
		t.Error("Expected to be able to add self-loop")
	}
	
	if !n.HasEdgeConnecting("A", "A") {
		t.Error("Expected self-loop to exist")
	}
	
	endpoints := n.IncidentNodes("self")
	if endpoints.NodeU != "A" || endpoints.NodeV != "A" {
		t.Error("Expected self-loop to connect A to A")
	}
}

func TestNetworkPanicOnSelfLoop(t *testing.T) {
	n := UndirectedNetwork[string, string]()
	n.AddNode("A")
	
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when adding self-loop to network that doesn't allow them")
		}
	}()
	
	n.AddEdge("self", "A", "A")
}

func TestNetworkWithParallelEdges(t *testing.T) {
	n := NewNetworkBuilder[string, string]().AllowParallelEdges().Build()
	
	n.AddNode("A")
	n.AddNode("B")
	
	// Add multiple edges between same nodes
	n.AddEdge("edge1", "A", "B")
	n.AddEdge("edge2", "A", "B")
	
	edges := n.EdgesConnecting("A", "B")
	if edges.Size() != 2 {
		t.Errorf("Expected 2 parallel edges, got %d", edges.Size())
	}
	
	if !edges.Contains("edge1") || !edges.Contains("edge2") {
		t.Error("Expected both edge1 and edge2 to connect A and B")
	}
}

func TestNetworkPanicOnParallelEdges(t *testing.T) {
	n := UndirectedNetwork[string, string]()
	n.AddNode("A")
	n.AddNode("B")
	n.AddEdge("edge1", "A", "B")
	
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when adding parallel edge to network that doesn't allow them")
		}
	}()
	
	n.AddEdge("edge2", "A", "B")
}

func TestNetworkAdjacentEdges(t *testing.T) {
	n := UndirectedNetwork[string, string]()
	
	n.AddNode("A")
	n.AddNode("B")
	n.AddNode("C")
	
	n.AddEdge("e1", "A", "B")
	n.AddEdge("e2", "B", "C")
	n.AddEdge("e3", "A", "C")
	
	// Test adjacent edges to e1 (should include e2 and e3)
	adjacentEdges := n.AdjacentEdges("e1")
	if adjacentEdges.Size() != 2 {
		t.Errorf("Expected 2 adjacent edges to e1, got %d", adjacentEdges.Size())
	}
	if !adjacentEdges.Contains("e2") || !adjacentEdges.Contains("e3") {
		t.Error("Expected e2 and e3 to be adjacent to e1")
	}
}

func TestNetworkClear(t *testing.T) {
	n := UndirectedNetwork[string, string]()
	
	n.AddNode("A")
	n.AddNode("B")
	n.AddEdge("edge1", "A", "B")
	
	n.Clear()
	
	if !n.IsEmpty() {
		t.Error("Expected network to be empty after clear")
	}
	
	if n.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", n.Size())
	}
	
	if n.Edges().Size() != 0 {
		t.Errorf("Expected 0 edges after clear, got %d", n.Edges().Size())
	}
}

func TestNetworkString(t *testing.T) {
	n := DirectedNetwork[string, string]()
	n.AddNode("A")
	n.AddNode("B")
	n.AddEdge("edge1", "A", "B")
	
	str := n.String()
	if str == "" {
		t.Error("Expected non-empty string representation")
	}
	
	// Should contain information about being directed
	if len(str) < 5 {
		t.Error("Expected meaningful string representation")
	}
}

func TestNetworkAsGraphView(t *testing.T) {
	network := UndirectedNetwork[string, string]()
	network.AddNode("A")
	network.AddNode("B")
	network.AddEdge("edge1", "A", "B")
	
	// Get graph view
	graph := network.AsGraph()
	
	// Test that graph view reflects the network
	if !graph.HasEdgeConnecting("A", "B") {
		t.Error("Expected graph view to have edge A-B")
	}
	
	if graph.Size() != 2 {
		t.Errorf("Expected graph view to have 2 nodes, got %d", graph.Size())
	}
	
	// Test that modifications through graph view affect network
	graph.AddNode("C")
	if !network.Contains("C") {
		t.Error("Expected network to contain node C added through graph view")
	}
	
	// Test that basic graph operations work
	successors := graph.Successors("A")
	if !successors.Contains("B") {
		t.Error("Expected B to be successor of A in graph view")
	}
	
	// Test that PutEdge panics (not supported in network graph view)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when calling PutEdge on network graph view")
		}
	}()
	
	graph.PutEdge("A", "C")
}