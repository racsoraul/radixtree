package radixtree

import (
	"testing"
)

func TestTree_SizeInconsistency(t *testing.T) {
	tree := New[int]()

	// 1. Set "apple"
	tree.Set("apple", 1)
	// root -> apple(1)
	if tree.root.size != 1 {
		t.Errorf("Expected root size 1 after apple, got %d", tree.root.size)
	}

	// 2. Set "apply" -> This should split "apple" into "appl" -> "e" and "y"
	tree.Set("apply", 2)
	// root -> appl(0) -> e(1), y(1)
	// "appl" is not a key yet, so its size should be 2.

	if len(tree.root.children) != 1 {
		t.Fatalf("Expected 1 child for root, got %d", len(tree.root.children))
	}

	applNode := tree.root.children[0].destination
	if applNode.size != 2 {
		t.Errorf("Expected appl node size 2 after apply, got %d", applNode.size)
	}
	if applNode.isKey {
		t.Errorf("Expected appl node NOT to be a key")
	}

	// 3. Set "appl" -> This is an exact match on "appl" node.
	// It should become a key.
	tree.Set("appl", 3)

	if !applNode.isKey {
		t.Errorf("Expected appl node to be a key after Set(\"appl\")")
	}

	// BUG: Expected appl node size 3 after becoming a key, but it was 2.
	if applNode.size != 3 {
		t.Errorf("Expected appl node size 3 after becoming a key, got %d", applNode.size)
	}

	if tree.root.size != 3 {
		t.Errorf("Expected root size 3, got %d", tree.root.size)
	}
}
