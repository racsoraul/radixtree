package radixtree

// Tree Represents a Radix Tree. Operations are not concurrency safe.
type Tree struct {
	root *Node
	size int
}

// New Creates a new Radix Tree.
func New() *Tree {
	return &Tree{
		root: NewNode(false, nil),
	}
}

// Size Returns the number of entries in the tree.
func (t *Tree) Size() int {
	return t.size
}

// String Returns a string representation of the Tree.
func (t *Tree) String() string {
	if t == nil {
		return ""
	}
	return t.root.string(0)
}

// Insert Adds a new entry to the tree or updates an existing one.
func (t *Tree) Insert(entry string, data any) {
	if entry == "" {
		t.root.isKey = true
		t.root.data = data
		return
	}

	currentNode := t.root

	for {
		edge, prefix, entrySuffix, edgeSuffix := currentNode.matchEdge(entry)
		if edge == nil {
			// No edge found. Create a new one.
			currentNode.children[entry[0]] = NewEdge(entry, NewNode(true, data))
			t.size++
			return
		}

		if entrySuffix == "" && edgeSuffix == "" {
			// Exact match.
			if !edge.destination.isKey {
				edge.destination.isKey = true
				t.size++
			}
			edge.destination.data = data
			return
		}

		if entrySuffix == "" {
			// The entry is a prefix of the label. New node before child.
			entryEdge := NewEdge(prefix, NewNode(true, data))
			currentNode.children[prefix[0]] = entryEdge
			edge.label = edgeSuffix
			entryEdge.destination.children[edgeSuffix[0]] = edge
			t.size++
			return
		}

		if edgeSuffix == "" {
			// Label is a prefix of the entry. Traverse the edge to the child node.
			currentNode = edge.destination
			entry = entrySuffix
			continue
		}

		// There's a common prefix. We need to split the edge.
		bridge := NewEdge(prefix, NewNode(false, nil))
		currentNode.children[prefix[0]] = bridge
		edge.label = edgeSuffix
		bridge.destination.children[edgeSuffix[0]] = edge
		newEntryNode := NewEdge(entrySuffix, NewNode(true, data))
		bridge.destination.children[entrySuffix[0]] = newEntryNode
		t.size++
		return
	}
}

// Search Returns the data and true if the entry is in the tree. Returns nil and false otherwise.
func (t *Tree) Search(entry string) (any, bool) {
	if entry == "" {
		return t.root.data, t.root.isKey
	}

	currentNode := t.root

	for {
		edge, _, entrySuffix, edgeSuffix := currentNode.matchEdge(entry)
		if edge == nil {
			return nil, false
		}

		if entrySuffix != "" && edgeSuffix != "" {
			// Partial match.
			return nil, false
		}

		if entrySuffix == "" {
			if edgeSuffix == "" {
				// Exact match.
				return edge.destination.data, edge.destination.isKey
			}
			// Partial match. The entry is not a key node.
			return nil, false
		}

		// Move to the next child node.
		currentNode = edge.destination
		entry = entrySuffix
	}
}

// LongestPrefix Returns the longest prefix of the entry that is also a key in the tree.
func (t *Tree) LongestPrefix(entry string) string {
	if entry == "" {
		return ""
	}

	currentNode := t.root
	longestPrefix := ""

	for {
		edge, prefix, entrySuffix, edgeSuffix := currentNode.matchEdge(entry)
		if edge == nil {
			return longestPrefix
		}

		if entrySuffix != "" && edgeSuffix != "" {
			// Partial match.
			return longestPrefix
		}

		if entrySuffix == "" {
			if edgeSuffix == "" {
				// Exact match.
				return longestPrefix + prefix
			}
			// Partial match. The entry is not a key node.
			return longestPrefix
		}

		// Move to the next child node.
		currentNode = edge.destination
		entry = entrySuffix
		longestPrefix += prefix
	}
}
