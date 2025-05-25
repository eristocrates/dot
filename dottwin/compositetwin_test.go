package dottwin

import (
	"os"
	"testing"

	"github.com/eristocrates/dot"
)

func TestExampleSubsystemTwinGraph(t *testing.T) {
	g := NewGraphTwin(dot.Directed)

	c1 := g.NodeTwin("component")

	twinsub := NewCompositeTwin("testing/", "subsystem", g)
	twinsub.InputTwin("in1", c1)
	twinsub.InputTwin("in2", c1)
	twinsub.OutputTwin("out2", c1)

	twinsub.Export(func(g *dot.Graph) {
		sc1 := twinsub.NodeTwin("subcomponent 1")
		sc2 := twinsub.NodeTwin("subcomponent 2")
		twinsub.InputTwin("in1", sc1)
		twinsub.OutputTwin("out2", sc2)
		sc1.Edge(sc2)

		sub2 := NewCompositeTwin("testing/", "subsystem2", &GraphTwin{in: twinsub.in.Graph, ex: twinsub.ex.Graph})
		sub2.Export(func(g *dot.Graph) {
			sub2.InputTwin("in3", sc1)
			sub2.OutputTwin("out3", sc2)
			sub3 := sub2.NodeTwin("subcomponent 3")
			sub2.InputTwin("in3", sub3)
		})
	})
	os.WriteFile("testing/TestExampleTwinSubsystemSameGraph.dot", []byte(g.in.String()), os.ModePerm)
	os.WriteFile("testing/TestExampleTwinSubsystemExternalGraph.dot", []byte(g.ex.String()), os.ModePerm)
	os.WriteFile("testing/TestExampleTwinSubsystemMermaidDiagram.mmd", []byte(dot.MermaidFlowchart(g.in, dot.MermaidTopDown)), os.ModePerm)
}
