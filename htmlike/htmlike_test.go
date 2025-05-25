package htmlike

import (
	"fmt"
	"os"
	"testing"

	"github.com/eristocrates/dot"
)

func TestExampleHtmLike(t *testing.T) {

	g := dot.NewGraph(dot.Directed)
	g.SetAttribute("rankdir", "LR")

	// Set default shape for nodes, which HTML-like labels often override visually
	g.NodeInitializer(
		func(n dot.Node) {
			n.SetAttribute("shape", "plaintext")
		})

	// --- Node a ---
	aNode := g.Node("a")
	aLabel := NewHtmLike(
		NewTable().
			SetBorder(0).      // BORDER="0"
			SetCellBorder(1).  // CELLBORDER="1"
			SetCellSpacing(0). // CELLSPACING="0"
			AppendChildren(
				// First Row
				NewTR().AppendChildren(
					NewTD().
						SetRowspan(3).                    // ROWSPAN="3"
						SetBGColor("yellow").             // BGCOLOR="yellow"
						AppendChildren(NewText("class")), // Content "class"
				),
				// Second Row
				NewTR().AppendChildren(
					NewTD().
						SetPort("here").                      // PORT="here"
						SetBGColor("lightblue").              // BGCOLOR="lightblue"
						AppendChildren(NewText("qualifier")), // Content "qualifier"
				),
				// Note: The third logical row implied by ROWSPAN="3" on the first TD
				// is not explicitly represented as a TR in the source DOT, which is
				// unusual for standard HTML but seems how Graphviz interprets it here.
				// The structure directly translates the source: two TRs, each with one TD.
			),
	).String() // Generate the <...> string
	aNode.SetAttribute("label", dot.HTML(aLabel)) // Set the node's label attribute

	// --- Node b ---
	bNode := g.Node("b")
	// Node b overrides the default shape and adds a style
	bNode.SetAttribute("shape", "ellipse")
	bNode.SetAttribute("style", "filled")
	bLabel := NewHtmLike(
		NewTable().
			SetBGColor("bisque"). // BGCOLOR="bisque"
			AppendChildren(
				// First Row
				NewTR().AppendChildren(
					// Cell 1: "elephant" spanning 3 columns
					NewTD().SetColspan(3).AppendChildren(NewText("elephant")), // COLSPAN="3"
					// Cell 2: "two" spanning 2 rows, specific styles
					NewTD().
						SetRowspan(2).                  // ROWSPAN="2"
						SetBGColor("chartreuse").       // BGCOLOR="chartreuse"
						SetValign(ValignBOTTOM).        // VALIGN="bottom"
						SetAlign(AlignRIGHT).           // ALIGN="right"
						AppendChildren(NewText("two")), // Content "two"
				),
				// Second Row
				NewTR().AppendChildren(
					// Cell 1: Nested table, spanning 2 columns and 2 rows
					NewTD().
						SetColspan(2). // COLSPAN="2"
						SetRowspan(2). // ROWSPAN="2"
						AppendChildren(
							// Nested Table
							NewTable().
								SetBGColor("grey"). // BGCOLOR="grey"
								AppendChildren(
									// Nested Row 1: "corn"
									NewTR().AppendChildren(NewTD().AppendChildren(NewText("corn"))),
									// Nested Row 2: "c" with yellow background
									NewTR().AppendChildren(NewTD().SetBGColor("yellow").AppendChildren(NewText("c"))), // BGCOLOR="yellow"
									// Nested Row 3: "f"
									NewTR().AppendChildren(NewTD().AppendChildren(NewText("f"))),
								),
						),
					// Cell 2: "penguin" with white background
					NewTD().SetBGColor("white").AppendChildren(NewText("penguin")), // BGCOLOR="white"
				),
				// Third Row (This row only has one cell defined, spanning columns 1+2)
				NewTR().AppendChildren(
					// Cell 1: "4" spanning 2 columns, specific styles and port
					NewTD().
						SetColspan(2).                // COLSPAN="2"
						SetBorder(4).                 // BORDER="4"
						SetAlign(AlignRIGHT).         // ALIGN="right"
						SetPort("there").             // PORT="there"
						AppendChildren(NewText("4")), // Content "4"
				),
			),
	).String()
	bNode.SetAttribute("label", dot.HTML(bLabel))

	// --- Node c ---
	cNode := g.Node("c")
	cLabel := NewHtmLike(
		NewText("long line 1"),       // First line of text
		NewBR(),                      // <BR/> - default align center
		NewText("line 2"),            // Second line of text
		NewBR().SetAlign(AlignLEFT),  // <BR ALIGN="LEFT"/>
		NewText("line 3"),            // Third line of text
		NewBR().SetAlign(AlignRIGHT), // <BR ALIGN="RIGHT"/>
	).String()
	cNode.SetAttribute("label", dot.HTML(cLabel))

	// --- Node d ---
	dNode := g.Node("d")
	dNode.SetAttribute("shape", "triangle")

	// --- Edges ---

	// Edge a -> b using ports
	// Assuming the dot package allows specifying ports on edge endpoints like this:
	edgeAB := g.EdgeWithPorts(aNode, bNode, "here", "there")
	edgeAB.SetAttribute("dir", "both")
	edgeAB.SetAttribute("arrowtail", "diamond")

	// Edge c -> b
	g.Edge(cNode, bNode)

	// Edge d -> c with HTML label on edge
	edgeDC := g.Edge(dNode, cNode)
	edgeLabelDC := NewHtmLike(
		NewTable().AppendChildren(
			NewTR().AppendChildren(
				// Left colored square cell
				NewTD().SetBGColor("red").SetWidth(10).AppendChildren(NewText(" ")), // Need content, space is used
				// Middle text cell with line break
				NewTD().AppendChildren(
					NewText("Edge labels"),
					NewBR(),
					NewText("also"),
				),
				// Right colored square cell
				NewTD().SetBGColor("blue").SetWidth(10).AppendChildren(NewText(" ")), // Need content, space is used
			),
		),
	).String()
	edgeDC.SetAttribute("label", dot.HTML(edgeLabelDC))

	// --- Subgraph ---
	// Group nodes b and c in a subgraph for ranking
	// Assuming your dot package has a method to create a subgraph off the main graph
	subgraphSameRank := g.Subgraph("") // Anonymous subgraph
	subgraphSameRank.AddToSameRank(bNode, cNode)
	// Add the existing nodes to the subgraph

	// Print the generated DOT code
	fmt.Fprintln(os.Stdout, g.String())
	os.WriteFile("testing/TestExampleHtmlike.dot", []byte(g.String()), os.ModePerm)

}
