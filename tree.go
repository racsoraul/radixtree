package radixtree

// Tree Represents a Radix Tree.
type Tree struct {
	root *Node
}

// New Creates a new Radix Tree.
func New() *Tree {
	return &Tree{
		root: NewNode(false),
	}
}

// String Returns a string representation of the Tree.
func (t *Tree) String() string {
	if t == nil {
		return ""
	}
	return t.root.string(0)
}

// Insert Adds new text to the tree.
func (t *Tree) Insert(text string) {
	if text == "" {
		t.root.isKey = true
		return
	}

	currentNode := t.root

	for {
		edge, prefix, textSuffix, edgeSuffix := currentNode.matchEdge(text)
		if edge == nil {
			// No edge found. Create a new one.
			currentNode.children[text[0]] = NewEdge(text, NewNode(true))
			return
		}

		if textSuffix == "" && edgeSuffix == "" {
			// Exact match.
			edge.destination.isKey = true
			return
		}

		if textSuffix == "" {
			// The text is a prefix of the label. New node before child.
			textEdge := NewEdge(prefix, NewNode(true))
			currentNode.children[prefix[0]] = textEdge
			edge.label = edgeSuffix
			textEdge.destination.children[edgeSuffix[0]] = edge
			return
		}

		if edgeSuffix == "" {
			// Label is a prefix of text. Traverse the edge to the child node.
			currentNode = edge.destination
			text = textSuffix
			continue
		}

		// There's a common prefix. We need to split the edge.
		bridge := NewEdge(prefix, NewNode(false))
		currentNode.children[prefix[0]] = bridge
		edge.label = edgeSuffix
		bridge.destination.children[edgeSuffix[0]] = edge
		newTextNode := NewEdge(textSuffix, NewNode(true))
		bridge.destination.children[textSuffix[0]] = newTextNode
		return
	}
}
