package radixtree

import "strings"

// Node Represents a node in the Radix Tree.
type Node struct {
	isKey    bool
	children map[byte]*Edge
}

// NewNode Creates a new Node.
func NewNode(isKey bool) *Node {
	return &Node{
		isKey:    isKey,
		children: make(map[byte]*Edge),
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

// matchEdge Returns the edge that matches the first character of text, if any,
// the longest common prefix of text and edge label, and the suffixes of text and edge label.
func (n *Node) matchEdge(text string) (matchedEdge *Edge, commonPrefix, textSuffix, edgeSuffix string) {
	edge, ok := n.children[text[0]]
	if !ok {
		return nil, "", text, ""
	}

	commonPrefix, textSuffix, edgeSuffix = longestCommonPrefix(text, edge.label)
	return edge, commonPrefix, textSuffix, edgeSuffix
}

// longestCommonPrefix Returns the longest common prefix of text and label.
// The returned suffixes are the suffixes of text and label that are not part of the prefix.
func longestCommonPrefix(text, label string) (prefix, textSuffix, edgeSuffix string) {
	minLen := min(len(text), len(label))

	for i := 0; i < minLen; i++ {
		if text[i] != label[i] {
			return text[:i], text[i:], label[i:]
		}
	}

	return text[:minLen], text[minLen:], label[minLen:]
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
		sb.WriteString(indent + label + "\n")
		sb.WriteString(edge.destination.string(nextIndent))
	}
	return sb.String()
}
