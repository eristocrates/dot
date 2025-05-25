package dottwin

import "github.com/eristocrates/dot"

type NodeTwin struct {
	in dot.Node
	ex dot.Node
}

func (n NodeTwin) SetAttribute(label string, value interface{}) NodeTwin {
	n.in.AttributesMap.SetAttribute(label, value)
	n.ex.AttributesMap.SetAttribute(label, value)
	return n
}

// Edge sets label=value and returns the Edge for chaining.
func (n NodeTwin) Edge(toNodeTwin NodeTwin, labels ...string) EdgeTwin {
	return EdgeTwin{in: n.in.Graph().Edge(n.in, toNodeTwin.in, labels...), ex: n.ex.Graph().Edge(n.ex, toNodeTwin.ex, labels...)}
}
