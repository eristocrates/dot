package dottwin

import "github.com/eristocrates/dot"

type GraphTwin struct {
	in *dot.Graph
	ex *dot.Graph
}

func NewGraphTwin(options ...dot.GraphOption) *GraphTwin {
	return &GraphTwin{
		in: dot.NewGraph(options...),
		ex: dot.NewGraph(options...),
	}
}
func (g *GraphTwin) NodeTwin(id string) NodeTwin {
	return NodeTwin{in: g.in.Node(id),
		ex: g.ex.Node(id)}
}
