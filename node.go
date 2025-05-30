package dot

// Node represents a dot Node.
type Node struct {
	AttributesMap
	graph *Graph
	id    string
	seq   int
}

// ID returns the assigned id to this node.
func (n Node) ID() string { return n.id }

// SetAttribute sets label=value and return the Node
func (n Node) SetAttribute(label string, value interface{}) Node {
	n.AttributesMap.SetAttribute(label, value)
	return n
}

// Label sets the attribute "label" to the given label
func (n Node) Label(label string) Node {
	return n.SetAttribute("label", label)
}

// Box sets the attribute "shape" to "box"
func (n Node) Box() Node {
	return n.SetAttribute("shape", "box")
}

// Edge sets label=value and returns the Edge for chaining.
func (n Node) Edge(toNode Node, labels ...string) Edge {
	return n.graph.Edge(n, toNode, labels...)
}

// EdgesTo returns all existing edges between this Node and the argument Node.
func (n Node) EdgesTo(toNode Node) []Edge {
	return n.graph.FindEdges(n, toNode)
}

// Attribute returns the value stored by a name. Returns nil if missing.
func (n Node) Attribute(name string) interface{} {
	return n.attributes[name]
}

// ReverseEdge sets label=value and returns the Edge for chaining.
func (n Node) ReverseEdge(fromNode Node, labels ...string) Edge {
	return n.graph.Edge(fromNode, n, labels...)
}

// BidirectionalEdge adds two edges, marks the first as invisible and the second with direction "both". Returns both edges.
func (n Node) BidirectionalEdge(toAndFromNode Node) []Edge {
	e1 := n.Edge(toAndFromNode).SetAttribute("style", "invis")
	e2 := toAndFromNode.Edge(n).SetAttribute("dir", "both")
	return []Edge{e1, e2}
}

// NewRecordBuilder returns a new recordBuilder for setting the attributes for a record-shaped node.
// Call Build() on the builder to set the label and shape.
func (n Node) NewRecordBuilder() *recordBuilder {
	return newRecordBuilder(n)
}

func (n Node) Graph() *Graph {
	return n.graph
}
