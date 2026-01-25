package radixtree

import (
	"fmt"
	"sort"
	"strings"
)

// edge Connects nodes.
type edge struct {
	destination *node
	label       string
}

// newEdge Creates a new Edge.
func newEdge(label string, dest *node) edge {
	return edge{
		destination: dest,
		label:       label,
	}
}

// node Represents a node in the Radix Tree.
type node struct {
	children []edge
	isKey    bool
	size     int // Number of keys in this subtree (including this node if isKey).
	data     any
}

// newNode Creates a new node.
func newNode(isKey bool, data any) *node {
	size := 0
	if isKey {
		size = 1
	}
	return &node{
		children: make([]edge, 0),
		isKey:    isKey,
		size:     size,
		data:     data,
	}
}

// allKeys Populates all keys prefixed by prefix that exist in the node into the
// provided keys slice.
func (n *node) allKeys(prefix []byte, keys *[]string) {
	if n.isKey {
		*keys = append(*keys, string(prefix))
	}

	for _, child := range n.children {
		prevLen := len(prefix)
		prefix = append(prefix, child.label...)
		child.destination.allKeys(prefix, keys)
		prefix = prefix[:prevLen] // Restore buffer.
	}
}

// addEdge Adds an edge to the subtree.
func (n *node) addEdge(e edge) {
	num := len(n.children)
	index := sort.Search(num, func(i int) bool {
		return n.children[i].label[0] >= e.label[0]
	})

	n.children = append(n.children, edge{})
	copy(n.children[index+1:], n.children[index:])
	n.children[index] = e
}

// updateEdge Updates an existing edge in the subtree.
func (n *node) updateEdge(e edge) {
	num := len(n.children)
	index := sort.Search(num, func(i int) bool {
		return n.children[i].label[0] >= e.label[0]
	})
	if index < num && n.children[index].label[0] == e.label[0] {
		n.children[index] = e
		return
	}
	panic(e.label + ": edge not found")
}

// matchEdge Returns the edge that matches the first character of the entry, if any,
// the longest common prefix of entry and edge label, and the suffixes of entry and edge label.
func (n *node) matchEdge(entry string) (matchedEdge edge, commonPrefix, entrySuffix, edgeSuffix string) {
	num := len(n.children)
	index := sort.Search(num, func(i int) bool {
		return n.children[i].label[0] >= entry[0]
	})

	if index < num && n.children[index].label[0] == entry[0] {
		commonPrefix, entrySuffix, edgeSuffix = longestCommonPrefix(entry, n.children[index].label)
		return n.children[index], commonPrefix, entrySuffix, edgeSuffix
	}

	return edge{}, "", entry, ""
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
func (n *node) String() string {
	if n == nil {
		return ""
	}
	return n.string(0)
}

// string Returns a string representation of the node.
func (n *node) string(spacing int) string {
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
