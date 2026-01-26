// Package radixtree A fast, efficient Radix Tree implementation in Go.
// Provides a lexicographically ordered iteration and multiple lookup methods.
// It leverages Go iterators for a more natural API.
package radixtree

import "iter"

// labelBufferSize The size of a reusable buffer used to build complete labels.
const labelBufferSize = 64

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

// All Returns an iterator over all entries in the tree.
// It provides a lexicographically ordered iteration.
func (t *Tree) All() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		accumulatedLabel := make([]byte, 0, labelBufferSize)
		t.root.push(accumulatedLabel, yield)
	}
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
	panic("not implemented yet. Sorry!")
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

// KeysWithPrefix Returns a list of entry's keys in the tree that start with the given prefix.
// The limit parameter controls the maximum number of keys to return. If the value of limit is -1,
// all keys are returned.
func (t *Tree) KeysWithPrefix(prefix string, limit int) []string {
	if prefix == "" {
		return nil
	}

	currentNode := t.root
	accumulatedPrefix := make([]byte, 0, len(prefix)+labelBufferSize)

	for {
		matchedEdge, _, entrySuffix, edgeSuffix := currentNode.matchEdge(prefix)
		if matchedEdge.destination == nil {
			// No match.
			return nil
		}

		if entrySuffix != "" && edgeSuffix != "" {
			// Partial match due to a common prefix.
			return nil
		}

		// Add the matched edge label as the base prefix.
		accumulatedPrefix = append(accumulatedPrefix, matchedEdge.label...)

		if entrySuffix == "" {
			// Either an exact match or partial match.
			allocSize := matchedEdge.destination.size
			if limit > 0 && allocSize > limit {
				allocSize = limit
			}
			keys := make([]string, 0, allocSize)
			matchedEdge.destination.allKeys(accumulatedPrefix, &keys, limit)
			return keys
		}

		// Move to the next child node.
		currentNode = matchedEdge.destination
		prefix = entrySuffix
	}
}
