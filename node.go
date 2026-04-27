package radixtree

import (
	"fmt"
	"sort"
	"strings"
)

// edge Connects nodes.
type edge[T any] struct {
	destination *node[T]
	label       string
}

// newEdge Creates a new Edge.
func newEdge[T any](label string, dest *node[T]) edge[T] {
	return edge[T]{
		destination: dest,
		label:       label,
	}
}

// node Represents a node in the Radix Tree.
type node[T any] struct {
	children []edge[T]
	isKey    bool
	size     int // Number of keys in this subtree (including this node if isKey).
	data     T
}

// newNode Creates a new node.
func newNode[T any](isKey bool, data T) *node[T] {
	size := 0
	if isKey {
		size = 1
	}
	return &node[T]{
		children: make([]edge[T], 0),
		isKey:    isKey,
		size:     size,
		data:     data,
	}
}

// push Pushes all entries to the yield function.
func (n *node[T]) push(label []byte, yield func(string, T) bool) bool {
	if n.isKey {
		if !yield(string(label), n.data) {
			return false
		}
	}

	for _, child := range n.children {
		prevLen := len(label)
		label = append(label, child.label...)
		if !child.destination.push(label, yield) {
			return false
		}
		label = label[:prevLen] // Restore buffer.
	}

	return true
}

// allKeys Populates all keys prefixed by prefix that exist in the node into the
// provided keys slice.
func (n *node[T]) allKeys(prefix []byte, keys *[]string, limit int) {
	if n.isKey {
		if limit > 0 && len(*keys) >= limit {
			return
		}
		*keys = append(*keys, string(prefix))
	}

	for _, child := range n.children {
		prevLen := len(prefix)
		prefix = append(prefix, child.label...)
		child.destination.allKeys(prefix, keys, limit)
		prefix = prefix[:prevLen] // Restore buffer.
	}
}

// addEdge Adds an edge to the subtree.
func (n *node[T]) addEdge(e edge[T]) {
	num := len(n.children)
	index := sort.Search(num, func(i int) bool {
		return n.children[i].label[0] >= e.label[0]
	})

	n.children = append(n.children, edge[T]{})
	copy(n.children[index+1:], n.children[index:])
	n.children[index] = e
}

// updateEdge Updates an existing edge in the subtree.
func (n *node[T]) updateEdge(e edge[T]) error {
	num := len(n.children)
	index := sort.Search(num, func(i int) bool {
		return n.children[i].label[0] >= e.label[0]
	})
	if index < num && n.children[index].label[0] == e.label[0] {
		n.children[index] = e
		return nil
	}
	return fmt.Errorf("%s: edge not found", e.label)
}

// matchEdge Returns the edge that matches the first character of the entry, if any,
// the longest common prefix of entry and edge label, and the suffixes of entry and edge label.
func (n *node[T]) matchEdge(entry string) (matchedEdge edge[T], commonPrefix, entrySuffix, edgeSuffix string) {
	num := len(n.children)
	index := sort.Search(num, func(i int) bool {
		return n.children[i].label[0] >= entry[0]
	})

	if index < num && n.children[index].label[0] == entry[0] {
		commonPrefix, entrySuffix, edgeSuffix = longestCommonPrefix(entry, n.children[index].label)
		return n.children[index], commonPrefix, entrySuffix, edgeSuffix
	}

	return edge[T]{}, "", entry, ""
}

// longestCommonPrefix Returns the longest common prefix between entry and label.
// The returned suffixes are the suffixes of entry and label that are not part of the prefix.
func longestCommonPrefix(entry, label string) (prefix, entrySuffix, edgeSuffix string) {
	minLen := min(len(entry), len(label))

	for i := 0; i < minLen; i++ {
		if entry[i] != label[i] {
			return entry[:i], entry[i:], label[i:]
		}
	}

	return entry[:minLen], entry[minLen:], label[minLen:]
}

// String Returns a string representation of the node. It colors in
// yellow the key nodes.
func (n *node[T]) String() string {
	if n == nil {
		return ""
	}
	return n.string(0)
}

// string Returns a string representation of the node.
func (n *node[T]) string(spacing int) string {
	var sb strings.Builder
	var indent string
	if spacing > 0 {
		indent = strings.Repeat(" ", spacing-1)
		indent += "|__"
	}

	for _, child := range n.children {
		nextIndent := len(indent) + len(child.label)
		label := child.label
		if child.destination.isKey {
			label = "\x1b[33m" + label + "\x1b[0m"
		}
		sb.WriteString(indent + label + fmt.Sprintf("(%d)", child.destination.size) + "\n")
		sb.WriteString(child.destination.string(nextIndent))
	}
	return sb.String()
}
