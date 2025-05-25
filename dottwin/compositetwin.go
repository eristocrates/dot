package dottwin

import (
	"log"

	"github.com/eristocrates/dot"
	"github.com/eristocrates/dot/dotx"
)

type CompositeTwin struct {
	in dotx.Composite
	ex dotx.Composite
}

func NewCompositeTwin(dotPath, id string, g *GraphTwin) *CompositeTwin {
	// h := *g
	twinsub := &CompositeTwin{
		in: *dotx.NewComposite(dotPath, id, g.in, dotx.SameGraph),
		ex: *dotx.NewComposite(dotPath, id, g.ex, dotx.ExternalGraph),
	}
	return twinsub

}
func (c *CompositeTwin) Input(id string, from dotx.Connectable) dot.Edge {
	return c.in.Input(id, from)
}
func (c *CompositeTwin) Output(id string, to dotx.Connectable) dot.Edge {
	return c.in.Output(id, to)
}

/*
	func (c *CompositeTwin) Input(id string, from dotx.Connectable) EdgeTwin {
		return EdgeTwin{in: c.in.Input(id, from), ex: c.ex.Input(id, from)}
	}

	func (c *CompositeTwin) Output(id string, to dotx.Connectable) (inEdge dot.Edge, outEdge dot.Edge) {
		return c.in.Output(id, to), c.ex.Output(id, to)
	}
*/
func (c *CompositeTwin) InputTwin(id string, from NodeTwin) EdgeTwin {
	return EdgeTwin{in: c.in.Input(id, from.in), ex: c.ex.Input(id, from.ex)}
}
func (c *CompositeTwin) OutputTwin(id string, from NodeTwin) EdgeTwin {
	return EdgeTwin{in: c.in.Output(id, from.in), ex: c.ex.Output(id, from.ex)}
}
func (c *CompositeTwin) Export(build func(g *dot.Graph)) *CompositeTwin {
	build(c.ex.Graph)
	if err := c.ex.ExportFile(); err != nil {
		log.Println("WARN: dotx.Composite.Export failed", err)
	}
	return c
}

func (c *CompositeTwin) NodeTwin(id string) NodeTwin {
	return NodeTwin{in: c.in.Node(id),
		ex: c.ex.Node(id)}
}

func (c *CompositeTwin) SetAttribute(label string, value interface{}) NodeTwin {
	return NodeTwin{in: c.in.OuterNode().SetAttribute(label, value),
		ex: c.ex.OuterNode().SetAttribute(label, value)}
}
