package radixtree

import (
	"fmt"
	"strings"
)

// Node Represents a node in the Radix Tree.
type Node struct {
	children map[byte]*Edge
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
		children: make(map[byte]*Edge),
		isKey:    isKey,
		size:     size,
		data:     data,
	}
}

// Edge Connects nodes.
type Edge struct {
	destination *Node
	label       string
}

// NewEdge Creates a new Edge.
func NewEdge(label string, dest *Node) *Edge {
	return &Edge{
		destination: dest,
		label:       label,
	}
}

// allKeys Populates all keys prefixed by prefix that exist in the node into the
// provided keys slice.
func (n *Node) allKeys(prefix string, keys *[]string) {
	if n.isKey {
		*keys = append(*keys, prefix)
	}

	for _, edge := range n.children {
		edge.destination.allKeys(prefix+edge.label, keys)
	}
}

// matchEdge Returns the edge that matches the first character of the entry, if any,
// the longest common prefix of entry and edge label, and the suffixes of entry and edge label.
func (n *Node) matchEdge(entry string) (matchedEdge *Edge, commonPrefix, entrySuffix, edgeSuffix string) {
	edge, ok := n.children[entry[0]]
	if !ok {
		return nil, "", entry, ""
	}

	commonPrefix, entrySuffix, edgeSuffix = longestCommonPrefix(entry, edge.label)
	return edge, commonPrefix, entrySuffix, edgeSuffix
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
