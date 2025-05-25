package dotx

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/eristocrates/dot"
)

type compositeGraphKind int

// Connectable is a dot.Node or a *dotx.Composite
type Connectable interface {
	SetAttribute(label string, value interface{}) dot.Node
}

const (
	// SameGraph means that the composite graph will be a cluster within the graph.
	SameGraph compositeGraphKind = iota
	// ExternalGraph means the composite graph will be exported on its own, linked by the node within the graph
	ExternalGraph
)

// Composite is a graph and node to create abstractions in graphs.
type Composite struct {
	*dot.Graph
	outerNode   dot.Node
	outerGraph  *dot.Graph
	dotFilename string
	kind        compositeGraphKind
}

// NewComposite creates a Composite abstraction that is represented as a Node (box3d shape) in the graph.
// The kind determines whether the graph of the composite is embedded (same graph) or external.
func NewComposite(dotPath, id string, g *dot.Graph, kind compositeGraphKind) *Composite {
	var innerGraph *dot.Graph
	if kind == SameGraph {
		innerGraph = g.Subgraph(id, dot.ClusterOption{})
	} else {
		innerGraph = dot.NewGraph(dot.Directed)
	}
	sub := &Composite{
		Graph:      innerGraph,
		outerNode:  g.Node(id).SetAttribute("shape", "box3d"),
		outerGraph: g,
		kind:       kind,
	}
	sub.ExportName(dotPath, id)
	return sub
}

// ExportFilename returns the name of the file used by ExportFile. Override it using ExportName.
func (s *Composite) ExportFilename() string {
	return s.dotFilename
}

// SetAttribute sets label=value and returns the Node in the graph
func (s *Composite) SetAttribute(label string, value interface{}) dot.Node {
	return s.outerNode.SetAttribute(label, value)
}

// ExportName argument name will be used for the .dot export and the HREF link using svg
// So if name = "my example" then export will create "my_example.dot" and the link will be "my_example.svg"
func (s *Composite) ExportName(dotPath, name string) {
	hrefFile := strings.ReplaceAll(name, " ", "_") + ".svg"
	dotFile := strings.ReplaceAll(name, " ", "_") + ".dot"
	s.outerNode.SetAttribute("href", hrefFile)
	s.dotFilename = dotPath + dotFile
}

// Input creates an edge.
// If the from Connectable is part of the parent graph then the edge is added to the parent graph.
// If the from Connectable is part of the composite then the edge is added to the inner graph.
func (s *Composite) Input(id string, from Connectable) dot.Edge {
	var fromNode dot.Node
	if n, ok := from.(dot.Node); ok {
		fromNode = n
	} else {
		if c, ok := from.(*Composite); ok {
			fromNode = c.outerNode
		}
	}
	if s.Graph.HasNode(fromNode) {
		// edge on innergraph
		return s.connect(id, true, fromNode)
	}
	// ensure input node in innergraph
	s.Node(id).SetAttribute("shape", "point")
	// edge on outergraph
	return fromNode.Edge(s.outerNode).Label(id)
}

// Output creates an edge.
// If the to Connectable is part of the parent graph then the edge is added to the parent graph.
// If the to Connectable is part of the composite then the edge is added to the inner graph.
func (s *Composite) Output(id string, to Connectable) dot.Edge {
	var toNode dot.Node
	if n, ok := to.(dot.Node); ok {
		toNode = n
	} else {
		if c, ok := to.(*Composite); ok {
			toNode = c.outerNode
		}
	}
	if s.Graph.HasNode(toNode) {
		// edge on innergraph
		return s.connect(id, false, toNode)
	}
	// ensure output node in innergraph
	s.Node(id).SetAttribute("shape", "point")
	// edge on outergraph
	return s.outerNode.Edge(toNode).Label(id)
}

func (s *Composite) connect(portName string, isInput bool, inner dot.Node) dot.Edge {
	// node creation is idempotent
	port := s.Node(portName).SetAttribute("shape", "point")
	if isInput {
		return s.EdgeWithPorts(port, inner, "s", "n").SetAttribute("taillabel", portName)
	} else {
		// is output
		return s.EdgeWithPorts(inner, port, "s", "n").SetAttribute("headlabel", portName)
	}
}

// ExportFile creates a DOT file using the default name (based on name) or overridden using ExportName().
func (s *Composite) ExportFile() error {
	if s.kind != ExternalGraph {
		return errors.New("ExportFile is only applicable to a ExternalGraph Composite")
	}

	path := s.ExportFilename()
	dirPath, _ := filepath.Split(path)
	err := os.MkdirAll(dirPath, 0755)

	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(s.Graph.String()), os.ModePerm)
}

// Export writes the DOT file for a Composite after building the content (child) graph using the build function.
// Use ExportName() on the Composite to modify the filename used.
// If writing of the file fails then a warning is logged.
func (s *Composite) Export(build func(g *dot.Graph)) *Composite {
	build(s.Graph)
	if err := s.ExportFile(); err != nil {
		log.Println("WARN: dotx.Composite.Export failed", err)
	}
	return s
}

func (c *Composite) OuterNode() *dot.Node {
	return &c.outerNode
}

// ConvertExternalToSameGraph takes a root composite and converts it to same-graph format
func ConvertExternalToSameGraph(root *Composite) *dot.Graph {
	// Create new root graph
	newGraph := dot.NewGraph(dot.Directed)

	// Create a mapping of original nodes to new nodes
	nodeMap := make(map[string]dot.Node)

	// Copy all nodes from the root graph first
	for _, node := range root.outerGraph.FindNodes() {
		newNode := newGraph.Node(node.ID())
		copyAttributes(node, newNode)
		nodeMap[node.ID()] = newNode
	}

	// Now find and copy all edges by checking all node pairs
	nodes := root.outerGraph.FindNodes()
	for i, fromNode := range nodes {
		for _, toNode := range nodes[i+1:] {
			// Check edges in both directions
			edges := root.outerGraph.FindEdges(fromNode, toNode)
			for _, edge := range edges {
				newEdge := newGraph.Edge(nodeMap[edge.From().ID()], nodeMap[edge.To().ID()])
				copyAttributes(edge, newEdge)
			}

			// Check reverse direction
			edges = root.outerGraph.FindEdges(toNode, fromNode)
			for _, edge := range edges {
				newEdge := newGraph.Edge(nodeMap[edge.From().ID()], nodeMap[edge.To().ID()])
				copyAttributes(edge, newEdge)
			}
		}
	}

	// Recursively process composites
	var processComposite func(comp *Composite, parentGraph *dot.Graph, parentNodeMap map[string]dot.Node)
	processComposite = func(comp *Composite, parentGraph *dot.Graph, parentNodeMap map[string]dot.Node) {
		// Create cluster for this composite
		cluster := parentGraph.Subgraph(comp.outerNode.ID(), dot.ClusterOption{})
		cluster.SetAttribute("label", comp.outerNode.ID())
		cluster.SetAttribute("fillcolor", "white")

		// Create new node map for this cluster
		clusterNodeMap := make(map[string]dot.Node)

		// Copy all nodes from composite's inner graph
		for _, node := range comp.Graph.FindNodes() {
			newNode := cluster.Node(node.ID())
			copyAttributes(node, newNode)
			clusterNodeMap[node.ID()] = newNode
		}

		// Copy all edges from composite's inner graph
		compNodes := comp.Graph.FindNodes()
		for i, fromNode := range compNodes {
			for _, toNode := range compNodes[i+1:] {
				// Check edges in both directions
				edges := comp.Graph.FindEdges(fromNode, toNode)
				for _, edge := range edges {
					newEdge := cluster.Edge(clusterNodeMap[edge.From().ID()], clusterNodeMap[edge.To().ID()])
					copyAttributes(edge, newEdge)
				}

				// Check reverse direction
				edges = comp.Graph.FindEdges(toNode, fromNode)
				for _, edge := range edges {
					newEdge := cluster.Edge(clusterNodeMap[edge.From().ID()], clusterNodeMap[edge.To().ID()])
					copyAttributes(edge, newEdge)
				}
			}
		}

		// Process child composites
		for _, child := range findChildComposites(comp) {
			processComposite(child, cluster, clusterNodeMap)
		}
	}

	// Start processing from root composite
	processComposite(root, newGraph, nodeMap)

	return newGraph
}

// Helper function to copy attributes
func copyAttributes(from interface{}, to interface{}) {
	switch src := from.(type) {
	case dot.Node:
		if dst, ok := to.(dot.Node); ok {
			for key, value := range src.Attributes() {
				dst.SetAttribute(key, value)
			}
		}
	case dot.Edge:
		if dst, ok := to.(dot.Edge); ok {
			for key, value := range src.Attributes() {
				dst.SetAttribute(key, value)
			}
		}
	}
}

func findChildComposites(comp *Composite) []*Composite {
	var children []*Composite

	// Iterate over all nodes in the composite's inner graph
	for _, node := range comp.Graph.FindNodes() {
		attrs := node.Attributes()
		href, ok := attrs["href"].(string)
		if !ok || href == "" {
			continue
		}

		// Derive the .dot filename from the href (which is an SVG file)
		// e.g. "child_composite.svg" -> "child_composite.dot"
		// dotFile := strings.TrimSuffix(href, ".svg") + ".dot"
		dotPath := comp.dotFilename // full path of current composite's dot file
		// dir := filepath.Dir(dotPath)
		// fullDotPath := filepath.Join(dir, dotFile)

		// Create a new Composite for the child
		// Use ExternalGraph kind because child composites are external graphs linked by href
		child := NewComposite(dotPath, node.ID(), comp.Graph, SameGraph)

		// Optionally, you could load or parse the child's graph here if needed
		// but for now just create the Composite wrapper

		children = append(children, child)
	}

	return children
}
