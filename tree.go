package radixtree

// Tree Represents a Radix Tree. Operations are not concurrency safe.
type Tree struct {
	root *node
	size int
}

// New Creates a new Radix Tree.
func New() *Tree {
	return &Tree{
		root: newNode(false, nil),
	}
}

// Len Returns the number of entries (key nodes) in the tree.
func (t *Tree) Len() int {
	return t.size
}

// String Returns a string representation of the Tree.
func (t *Tree) String() string {
	if t == nil {
		return ""
	}
	return t.root.string(0)
}

// Set Adds a new entry to the tree or updates an existing one.
func (t *Tree) Set(entry string, data any) {
	if entry == "" {
		if !t.root.isKey {
			t.root.isKey = true
			t.root.size++
			t.size++
		}
		t.root.data = data
		return
	}

	path := make([]*node, 0, 1)
	currentNode := t.root

	for {
		path = append(path, currentNode) // Keep track of the path to update the sizes.
		matchedEdge, prefix, entrySuffix, edgeSuffix := currentNode.matchEdge(entry)
		if matchedEdge.destination == nil {
			// No edge found. Create a new one.
			currentNode.addEdge(newEdge(entry, newNode(true, data)))
			t.size++
			// Update sizes for all nodes in the path.
			for _, n := range path {
				n.size++
			}
			return
		}

		if entrySuffix == "" && edgeSuffix == "" {
			// Exact match.
			if !matchedEdge.destination.isKey {
				matchedEdge.destination.isKey = true
				t.size++
				// Update sizes for all nodes in the path.
				for _, n := range path {
					n.size++
				}
			}
			matchedEdge.destination.data = data
			return
		}

		if entrySuffix == "" {
			// The entry is a prefix of the label. New node before child.
			entryEdge := newEdge(prefix, newNode(true, data))
			currentNode.updateEdge(entryEdge)
			matchedEdge.label = edgeSuffix
			entryEdge.destination.addEdge(matchedEdge)
			t.size++

			// It gets the old subtree size plus itself (1).
			entryEdge.destination.size += matchedEdge.destination.size

			// Update sizes for all nodes in the path.
			for _, n := range path {
				n.size++
			}
			return
		}

		if edgeSuffix == "" {
			// Label is a prefix of the entry. Traverse the edge to the child node.
			currentNode = matchedEdge.destination
			entry = entrySuffix
			continue
		}

		// There's a common prefix. We need to split the edge.
		bridge := newEdge(prefix, newNode(false, nil))
		currentNode.updateEdge(bridge)
		matchedEdge.label = edgeSuffix
		bridge.destination.addEdge(matchedEdge)
		newEntryNode := newEdge(entrySuffix, newNode(true, data))
		bridge.destination.addEdge(newEntryNode)
		t.size++

		// Existing subtree size plus the new one (1).
		bridge.destination.size = matchedEdge.destination.size + 1

		// Update sizes for all nodes in the path.
		for _, n := range path {
			n.size++
		}
		return
	}
}

// Get Returns the data and true if the entry is in the tree. Returns nil and false otherwise.
func (t *Tree) Get(entry string) (any, bool) {
	if entry == "" {
		return t.root.data, t.root.isKey
	}

	currentNode := t.root

	for {
		matchedEdge, _, entrySuffix, edgeSuffix := currentNode.matchEdge(entry)
		if matchedEdge.destination == nil {
			return nil, false
		}

		if entrySuffix != "" && edgeSuffix != "" {
			// Partial match.
			return nil, false
		}

		if entrySuffix == "" {
			if edgeSuffix == "" {
				// Exact match.
				return matchedEdge.destination.data, matchedEdge.destination.isKey
			}
			// Partial match. The entry is not a key node.
			return nil, false
		}

		// Move to the next child node.
		currentNode = matchedEdge.destination
		entry = entrySuffix
	}
}

// Delete Removes an entry from the tree. Returns true indicating
// if the entry existed and was deleted, false otherwise.
func (t *Tree) Delete(entry string) bool {
	// TODO: Implement.
	return false
}

// LongestPrefix Returns the longest prefix of the entry that is also a key node in the tree.
func (t *Tree) LongestPrefix(entry string) string {
	if entry == "" {
		return ""
	}

	currentNode := t.root
	longestPrefix := ""

	for {
		matchedEdge, prefix, entrySuffix, edgeSuffix := currentNode.matchEdge(entry)
		if matchedEdge.destination == nil {
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
		currentNode = matchedEdge.destination
		entry = entrySuffix
		longestPrefix += prefix
	}
}

// prefixBufferSize The size of the buffer used to build the prefix of keys returned by KeysWithPrefix.
const prefixBufferSize = 64

// KeysWithPrefix Returns a list of entry's keys in the tree that start with the given prefix.
func (t *Tree) KeysWithPrefix(prefix string) []string {
	if prefix == "" {
		return nil
	}

	currentNode := t.root
	accumulatedPrefix := make([]byte, 0, len(prefix)+prefixBufferSize)
	accumulatedPrefix = append(accumulatedPrefix, prefix...)

	for {
		matchedEdge, _, entrySuffix, edgeSuffix := currentNode.matchEdge(prefix)
		if matchedEdge.destination == nil {
			// No match.
			return nil
		}

		if entrySuffix != "" && edgeSuffix != "" {
			// Partial match.
			return nil
		}

		if entrySuffix == "" {
			// Allocate all keys at once.
			keys := make([]string, 0, matchedEdge.destination.size)
			if edgeSuffix == "" {
				// Exact match.
				matchedEdge.destination.allKeys(accumulatedPrefix, &keys)
				return keys
			}

			// Partial match. The entry is not a key node.
			accumulatedPrefix = append(accumulatedPrefix, edgeSuffix...)
			matchedEdge.destination.allKeys(accumulatedPrefix, &keys)
			return keys
		}

		// Move to the next child node.
		currentNode = matchedEdge.destination
		prefix = entrySuffix
		accumulatedPrefix = append(accumulatedPrefix, prefix...)
	}
}
