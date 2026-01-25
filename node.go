package radixtree

import (
	"fmt"
	"sort"
	"strings"
)

// Edge Connects nodes.
type Edge struct {
	destination *Node
	label       string
}

// NewEdge Creates a new Edge.
func NewEdge(label string, dest *Node) Edge {
	return Edge{
		destination: dest,
		label:       label,
	}
}

// Node Represents a node in the Radix Tree.
type Node struct {
	children []Edge
	isKey    bool
	size     int // Number of keys in this subtree (including this node if isKey).
	data     any
}

// NewNode Creates a new Node.
func NewNode(isKey bool, data any) *Node {
	size := 0
	if isKey {
		size = 1
	}
	return &Node{
		children: make([]Edge, 0),
		isKey:    isKey,
		size:     size,
		data:     data,
	}
}

// allKeys Populates all keys prefixed by prefix that exist in the node into the
// provided keys slice.
func (n *Node) allKeys(prefix []byte, keys *[]string) {
	if n.isKey {
		*keys = append(*keys, string(prefix))
	}

	for _, edge := range n.children {
		prevLen := len(prefix)
		prefix = append(prefix, edge.label...)
		edge.destination.allKeys(prefix, keys)
		prefix = prefix[:prevLen] // Restore buffer.
	}
}

// addEdge Adds an edge to the subtree.
func (n *Node) addEdge(edge Edge) {
	num := len(n.children)
	index := sort.Search(num, func(i int) bool {
		return n.children[i].label[0] >= edge.label[0]
	})

	n.children = append(n.children, Edge{})
	copy(n.children[index+1:], n.children[index:])
	n.children[index] = edge
}

// updateEdge Updates an existing edge in the subtree.
func (n *Node) updateEdge(edge Edge) {
	num := len(n.children)
	index := sort.Search(num, func(i int) bool {
		return n.children[i].label[0] >= edge.label[0]
	})
	if index < num && n.children[index].label[0] == edge.label[0] {
		n.children[index] = edge
		return
	}
	panic(edge.label + ": edge not found")
}

// matchEdge Returns the edge that matches the first character of the entry, if any,
// the longest common prefix of entry and edge label, and the suffixes of entry and edge label.
func (n *Node) matchEdge(entry string) (matchedEdge Edge, commonPrefix, entrySuffix, edgeSuffix string) {
	num := len(n.children)
	index := sort.Search(num, func(i int) bool {
		return n.children[i].label[0] >= entry[0]
	})

	if index < num && n.children[index].label[0] == entry[0] {
		commonPrefix, entrySuffix, edgeSuffix = longestCommonPrefix(entry, n.children[index].label)
		return n.children[index], commonPrefix, entrySuffix, edgeSuffix
	}

	return Edge{}, "", entry, ""
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

// String Returns a string representation of the Node. It colors in
// yellow the key nodes.
func (n *Node) String() string {
	if n == nil {
		return ""
	}
	return n.string(0)
}

// string Returns a string representation of the node.
func (n *Node) string(spacing int) string {
	var sb strings.Builder
	var indent string
	if spacing > 0 {
		indent = strings.Repeat(" ", spacing-1)
		indent += "|__"
	}

	for _, edge := range n.children {
		nextIndent := len(indent) + len(edge.label)
		label := edge.label
		if edge.destination.isKey {
			label = "\x1b[33m" + label + "\x1b[0m"
		}
		sb.WriteString(indent + label + fmt.Sprintf("(%d)", edge.destination.size) + "\n")
		sb.WriteString(edge.destination.string(nextIndent))
	}
	return sb.String()
}
