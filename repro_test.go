package radixtree

import (
	"testing"
)

func TestUpdateEdgeError(t *testing.T) {
	n := newNode[int](false, 0)
	// Try to update an edge that doesn't exist
	err := n.updateEdge(newEdge("nonexistent", newNode[int](true, 1)))
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
	expectedMsg := "nonexistent: edge not found"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message %q, got %q", expectedMsg, err.Error())
	}
}

func TestTreeSetError(t *testing.T) {
	// Although Set should not normally return an error if the tree is used through its public API,
	// we want to ensure it propagates errors from updateEdge if they were to happen.
	// Since updateEdge is not directly reachable with an error condition through Set (it's a bug if it is),
	// this is more of a structural check.
	tree := New[int]()
	err := tree.Set("hello", 1)
	if err != nil {
		t.Errorf("Expected no error from Set, got %v", err)
	}
}
