package htmlike

import (
	"bytes"
	"fmt"
	"html"
	"io"
)

// HtmLike represents a Graphviz HTML-like label, enclosed in <...>.
// It can contain one or more HtmLikeElement children.
// Typically, the root element is a *Table.
type HtmLike struct {
	Elements []HtmLikeElement
}

// HtmLikeElement is an interface for any element that can be part of an HTML-like label.
type HtmLikeElement interface {
	WriteDOT(w io.Writer) error
}

// Text represents simple text content within an HTML-like label.
// The text content needs to be HTML-escaped.
type Text string

// WriteDOT writes the escaped text content to the writer.
func (t Text) WriteDOT(w io.Writer) error {
	// Escape characters required by Graphviz HTML-like labels
	escaped := html.EscapeString(string(t))
	_, err := io.WriteString(w, escaped)
	return err
}

// String converts Text to its string representation.
func (t Text) String() string {
	return string(t)
}

// NewText creates a new Text element.
func NewText(s string) Text {
	return Text(s)
}

// Attribute value types and constants for specific attributes

type HtmLikeAlign string

const (
	AlignCENTER HtmLikeAlign = "CENTER"
	AlignLEFT   HtmLikeAlign = "LEFT"
	AlignRIGHT  HtmLikeAlign = "RIGHT"
	AlignTEXT   HtmLikeAlign = "TEXT" // TD only
)

type HtmLikeValign string

const (
	ValignMIDDLE HtmLikeValign = "MIDDLE"
	ValignBOTTOM HtmLikeValign = "BOTTOM"
	ValignTOP    HtmLikeValign = "TOP"
)

type HtmLikeImgScale string

const (
	ScaleFALSE  HtmLikeImgScale = "FALSE"
	ScaleTRUE   HtmLikeImgScale = "TRUE"
	ScaleWIDTH  HtmLikeImgScale = "WIDTH"
	ScaleHEIGHT HtmLikeImgScale = "HEIGHT"
	ScaleBOTH   HtmLikeImgScale = "BOTH"
)

// Helper function to write optional attributes
func writeAttrString(w io.Writer, name string, value string) error {
	if value == "" {
		return nil
	}
	// Attribute values must be in double quotes and escaped.
	// Graphviz spec says escString for HREF, ID, TARGET, TITLE, TOOLTIP.
	// Let's escape all string attribute values to be safe, consistent with HTML.
	escapedValue := html.EscapeString(value)
	_, err := fmt.Fprintf(w, ` %s="%s"`, name, escapedValue)
	return err
}

func writeAttrIntPtr(w io.Writer, name string, value *int) error {
	if value == nil {
		return nil
	}
	_, err := fmt.Fprintf(w, ` %s="%d"`, name, *value)
	return err
}

func writeAttrUintPtr(w io.Writer, name string, value *uint) error {
	if value == nil {
		return nil
	}
	_, err := fmt.Fprintf(w, ` %s="%d"`, name, *value)
	return err
}

func writeAttrBoolPtr(w io.Writer, name string, value *bool) error {
	if value == nil {
		return nil
	}
	valStr := "FALSE"
	if *value {
		valStr = "TRUE"
	}
	_, err := fmt.Fprintf(w, ` %s="%s"`, name, valStr)
	return err
}

func writeAttrHtmLikeAlign(w io.Writer, name string, value HtmLikeAlign) error {
	if value == "" {
		return nil
	}
	_, err := fmt.Fprintf(w, ` %s="%s"`, name, string(value))
	return err
}

func writeAttrHtmLikeValign(w io.Writer, name string, value HtmLikeValign) error {
	if value == "" {
		return nil
	}
	_, err := fmt.Fprintf(w, ` %s="%s"`, name, string(value))
	return err
}

func writeAttrHtmLikeImgScale(w io.Writer, name string, value HtmLikeImgScale) error {
	if value == "" {
		return nil
	}
	_, err := fmt.Fprintf(w, ` %s="%s"`, name, string(value))
	return err
}

// Table represents a <TABLE> element.
type Table struct {
	Align         HtmLikeAlign
	BGColor       string
	Border        *uint // Note: BORDER is uint up to 255 for TABLE, 127 for CELLBORDER
	CellBorder    *uint // Max 127
	CellPadding   *uint // Max 255
	CellSpacing   *uint // Max 127
	Color         string
	Columns       string // Only "*" currently
	FixedSize     *bool
	GradientAngle *int
	Height        *uint // Max 65535
	HREF          string
	ID            string
	Port          string
	Rows          string // Only "*" currently
	Sides         string // e.g., "LT", "B", "LRTB"
	Style         string // e.g., "ROUNDED", "RADIAL"
	Target        string
	Title         string // Alias for TOOLTIP
	Tooltip       string
	Valign        HtmLikeValign
	Width         *uint // Max 65535

	Children []HtmLikeElement // Expected to contain *TR elements
}

// NewTable creates a new Table element.
func NewTable(children ...HtmLikeElement) *Table {
	return &Table{Children: children}
}

// Setters for fluent API
func (t *Table) SetAlign(v HtmLikeAlign) *Table   { t.Align = v; return t }
func (t *Table) SetBGColor(v string) *Table       { t.BGColor = v; return t }
func (t *Table) SetBorder(v uint) *Table          { t.Border = &v; return t }
func (t *Table) SetCellBorder(v uint) *Table      { t.CellBorder = &v; return t }
func (t *Table) SetCellPadding(v uint) *Table     { t.CellPadding = &v; return t }
func (t *Table) SetCellSpacing(v uint) *Table     { t.CellSpacing = &v; return t }
func (t *Table) SetColor(v string) *Table         { t.Color = v; return t }
func (t *Table) SetColumns(v string) *Table       { t.Columns = v; return t } // Should be "*"
func (t *Table) SetFixedSize(v bool) *Table       { t.FixedSize = &v; return t }
func (t *Table) SetGradientAngle(v int) *Table    { t.GradientAngle = &v; return t }
func (t *Table) SetHeight(v uint) *Table          { t.Height = &v; return t }
func (t *Table) SetHREF(v string) *Table          { t.HREF = v; return t }
func (t *Table) SetID(v string) *Table            { t.ID = v; return t }
func (t *Table) SetPort(v string) *Table          { t.Port = v; return t }
func (t *Table) SetRows(v string) *Table          { t.Rows = v; return t } // Should be "*"
func (t *Table) SetSides(v string) *Table         { t.Sides = v; return t }
func (t *Table) SetStyle(v string) *Table         { t.Style = v; return t } // e.g. "ROUNDED", "RADIAL"
func (t *Table) SetTarget(v string) *Table        { t.Target = v; return t }
func (t *Table) SetTitle(v string) *Table         { t.Title = v; return t } // Alias for Tooltip
func (t *Table) SetTooltip(v string) *Table       { t.Tooltip = v; return t }
func (t *Table) SetValign(v HtmLikeValign) *Table { t.Valign = v; return t }
func (t *Table) SetWidth(v uint) *Table           { t.Width = &v; return t }

// AppendChildren adds children to the table.
func (t *Table) AppendChildren(children ...HtmLikeElement) *Table {
	t.Children = append(t.Children, children...)
	return t
}

// WriteDOT writes the <TABLE> element and its children to the writer.
func (t *Table) WriteDOT(w io.Writer) error {
	if _, err := io.WriteString(w, "\n<TABLE"); err != nil {
		return err
	}

	// Write attributes
	writeAttrHtmLikeAlign(w, "ALIGN", t.Align) // No error check to simplify, assuming underlying write handles it
	writeAttrString(w, "BGCOLOR", t.BGColor)
	writeAttrUintPtr(w, "BORDER", t.Border)
	writeAttrUintPtr(w, "CELLBORDER", t.CellBorder)
	writeAttrUintPtr(w, "CELLPADDING", t.CellPadding)
	writeAttrUintPtr(w, "CELLSPACING", t.CellSpacing)
	writeAttrString(w, "COLOR", t.Color)
	writeAttrString(w, "COLUMNS", t.Columns)
	writeAttrBoolPtr(w, "FIXEDSIZE", t.FixedSize)
	writeAttrIntPtr(w, "GRADIENTANGLE", t.GradientAngle)
	writeAttrUintPtr(w, "HEIGHT", t.Height)
	writeAttrString(w, "HREF", t.HREF)
	writeAttrString(w, "ID", t.ID)
	writeAttrString(w, "PORT", t.Port)
	writeAttrString(w, "ROWS", t.Rows)
	writeAttrString(w, "SIDES", t.Sides)
	writeAttrString(w, "STYLE", t.Style)
	writeAttrString(w, "TARGET", t.Target)
	// Title and Tooltip are aliases, write Tooltip if Title is empty
	tooltipVal := t.Tooltip
	if tooltipVal == "" && t.Title != "" {
		tooltipVal = t.Title
	}
	writeAttrString(w, "TOOLTIP", tooltipVal)
	writeAttrHtmLikeValign(w, "VALIGN", t.Valign)
	writeAttrUintPtr(w, "WIDTH", t.Width)

	if _, err := io.WriteString(w, ">"); err != nil {
		return err
	}

	// Write children
	for _, child := range t.Children {
		if err := child.WriteDOT(w); err != nil {
			return err
		}
	}

	_, err := io.WriteString(w, "\n</TABLE>")
	return err
}

// TR represents a <TR> element.
type TR struct {
	Children []HtmLikeElement // Expected to contain *TD elements or HR/VR
}

// NewTR creates a new TR element.
func NewTR(children ...HtmLikeElement) *TR {
	return &TR{Children: children}
}

// AppendChildren adds children to the table row.
func (tr *TR) AppendChildren(children ...HtmLikeElement) *TR {
	tr.Children = append(tr.Children, children...)
	return tr
}

// WriteDOT writes the <TR> element and its children to the writer.
func (tr *TR) WriteDOT(w io.Writer) error {
	if _, err := io.WriteString(w, "\n<TR>"); err != nil {
		return err
	}
	for _, child := range tr.Children {
		if err := child.WriteDOT(w); err != nil {
			return err
		}
	}
	_, err := io.WriteString(w, "</TR>")
	return err
}

// TD represents a <TD> element.
type TD struct {
	Align         HtmLikeAlign // Can also be TEXT
	BAlign        HtmLikeAlign // Default ALIGN for child BRs
	BGColor       string
	Border        *uint // Max 255
	CellPadding   *uint // Max 255
	CellSpacing   *uint // Max 127 (Note: TD doesn't have CellSpacing in spec? Maybe inherited?)
	Color         string
	Colspan       *uint // Max 65535
	FixedSize     *bool
	GradientAngle *int
	Height        *uint // Max 65535
	HREF          string
	ID            string
	Port          string
	Rowspan       *uint  // Max 65535
	Sides         string // e.g., "LT", "B", "LRTB"
	Style         string // Only "RADIAL" currently
	Target        string
	Title         string // Alias for TOOLTIP
	Tooltip       string
	Valign        HtmLikeValign
	Width         *uint // Max 65535

	Children []HtmLikeElement // Can contain anything a label can contain (Text, <TABLE>, <FONT>, <BR>, <IMG>, style tags)
}

// NewTD creates a new TD element.
func NewTD(children ...HtmLikeElement) *TD {
	return &TD{Children: children}
}

// Setters for fluent API
func (td *TD) SetAlign(v HtmLikeAlign) *TD  { td.Align = v; return td }
func (td *TD) SetBAlign(v HtmLikeAlign) *TD { td.BAlign = v; return td }
func (td *TD) SetBGColor(v string) *TD      { td.BGColor = v; return td }
func (td *TD) SetBorder(v uint) *TD         { td.Border = &v; return td }
func (td *TD) SetCellPadding(v uint) *TD    { td.CellPadding = &v; return td }

// SetCellSpacing seems incorrect based on spec, TD doesn't have it.
// func (td *TD) SetCellSpacing(v uint) *TD { td.CellSpacing = &v; return td }
func (td *TD) SetColor(v string) *TD         { td.Color = v; return td }
func (td *TD) SetColspan(v uint) *TD         { td.Colspan = &v; return td }
func (td *TD) SetFixedSize(v bool) *TD       { td.FixedSize = &v; return td }
func (td *TD) SetGradientAngle(v int) *TD    { td.GradientAngle = &v; return td }
func (td *TD) SetHeight(v uint) *TD          { td.Height = &v; return td }
func (td *TD) SetHREF(v string) *TD          { td.HREF = v; return td }
func (td *TD) SetID(v string) *TD            { td.ID = v; return td }
func (td *TD) SetPort(v string) *TD          { td.Port = v; return td }
func (td *TD) SetRowspan(v uint) *TD         { td.Rowspan = &v; return td }
func (td *TD) SetSides(v string) *TD         { td.Sides = v; return td }
func (td *TD) SetStyle(v string) *TD         { td.Style = v; return td } // e.g. "RADIAL"
func (td *TD) SetTarget(v string) *TD        { td.Target = v; return td }
func (td *TD) SetTitle(v string) *TD         { td.Title = v; return td } // Alias for Tooltip
func (td *TD) SetTooltip(v string) *TD       { td.Tooltip = v; return td }
func (td *TD) SetValign(v HtmLikeValign) *TD { td.Valign = v; return td }
func (td *TD) SetWidth(v uint) *TD           { td.Width = &v; return td }

// AppendChildren adds children to the table cell.
func (td *TD) AppendChildren(children ...HtmLikeElement) *TD {
	td.Children = append(td.Children, children...)
	return td
}

// WriteDOT writes the <TD> element and its children to the writer.
func (td *TD) WriteDOT(w io.Writer) error {
	if _, err := io.WriteString(w, "<TD"); err != nil {
		return err
	}

	// Write attributes
	writeAttrHtmLikeAlign(w, "ALIGN", td.Align)
	writeAttrHtmLikeAlign(w, "BALIGN", td.BAlign)
	writeAttrString(w, "BGCOLOR", td.BGColor)
	writeAttrUintPtr(w, "BORDER", td.Border)
	writeAttrUintPtr(w, "CELLPADDING", td.CellPadding)
	// writeAttrUintPtr(w, "CELLSPACING", td.CellSpacing) // Based on spec, TD doesn't have this attribute
	writeAttrString(w, "COLOR", td.Color)
	writeAttrUintPtr(w, "COLSPAN", td.Colspan)
	writeAttrBoolPtr(w, "FIXEDSIZE", td.FixedSize)
	writeAttrIntPtr(w, "GRADIENTANGLE", td.GradientAngle)
	writeAttrUintPtr(w, "HEIGHT", td.Height)
	writeAttrString(w, "HREF", td.HREF)
	writeAttrString(w, "ID", td.ID)
	writeAttrString(w, "PORT", td.Port)
	writeAttrUintPtr(w, "ROWSPAN", td.Rowspan)
	writeAttrString(w, "SIDES", td.Sides)
	writeAttrString(w, "STYLE", td.Style) // e.g. "RADIAL"
	writeAttrString(w, "TARGET", td.Target)
	// Title and Tooltip are aliases, write Tooltip if Title is empty
	tooltipVal := td.Tooltip
	if tooltipVal == "" && td.Title != "" {
		tooltipVal = td.Title
	}
	writeAttrString(w, "TOOLTIP", tooltipVal)
	writeAttrHtmLikeValign(w, "VALIGN", td.Valign)
	writeAttrUintPtr(w, "WIDTH", td.Width)

	if _, err := io.WriteString(w, ">"); err != nil {
		return err
	}

	// Write children
	for _, child := range td.Children {
		if err := child.WriteDOT(w); err != nil {
			return err
		}
	}

	_, err := io.WriteString(w, "</TD>")
	return err
}

// FONT represents a <FONT> element.
type FONT struct {
	Color     string
	Face      string
	PointSize *int // Max 65535 apparently for POINT-SIZE as int? Spec says "value". Let's assume int for now.

	Children []HtmLikeElement // Can contain Text, BR, IMG, style tags, other FONT, or even TABLES/TRs/TDs
}

// NewFONT creates a new FONT element.
func NewFONT(children ...HtmLikeElement) *FONT {
	return &FONT{Children: children}
}

// Setters for fluent API
func (f *FONT) SetColor(v string) *FONT  { f.Color = v; return f }
func (f *FONT) SetFace(v string) *FONT   { f.Face = v; return f }
func (f *FONT) SetPointSize(v int) *FONT { f.PointSize = &v; return f }

// AppendChildren adds children to the font element.
func (f *FONT) AppendChildren(children ...HtmLikeElement) *FONT {
	f.Children = append(f.Children, children...)
	return f
}

// WriteDOT writes the <FONT> element and its children to the writer.
func (f *FONT) WriteDOT(w io.Writer) error {
	if _, err := io.WriteString(w, "<FONT"); err != nil {
		return err
	}

	// Write attributes
	writeAttrString(w, "COLOR", f.Color)
	writeAttrString(w, "FACE", f.Face)
	writeAttrIntPtr(w, "POINT-SIZE", f.PointSize)

	if _, err := io.WriteString(w, ">"); err != nil {
		return err
	}

	// Write children
	for _, child := range f.Children {
		if err := child.WriteDOT(w); err != nil {
			return err
		}
	}

	_, err := io.WriteString(w, "</FONT>")
	return err
}

// BR represents a <BR/> element.
type BR struct {
	Align HtmLikeAlign
}

// NewBR creates a new BR element.
func NewBR(align ...HtmLikeAlign) *BR {
	b := &BR{}
	if len(align) > 0 {
		b.Align = align[0]
	}
	return b
}

// SetAlign sets the align attribute for the BR element.
func (b *BR) SetAlign(v HtmLikeAlign) *BR { b.Align = v; return b }

// WriteDOT writes the <BR/> element to the writer.
func (b *BR) WriteDOT(w io.Writer) error {
	if _, err := io.WriteString(w, "<BR"); err != nil {
		return err
	}
	writeAttrHtmLikeAlign(w, "ALIGN", b.Align)
	_, err := io.WriteString(w, "/>") // Self-closing
	return err
}

// IMG represents an <IMG/> element.
type IMG struct {
	Scale HtmLikeImgScale
	Src   string
}

// NewIMG creates a new IMG element.
func NewIMG(src string) *IMG {
	return &IMG{Src: src}
}

// Setters for fluent API
func (img *IMG) SetScale(v HtmLikeImgScale) *IMG { img.Scale = v; return img }
func (img *IMG) SetSrc(v string) *IMG            { img.Src = v; return img }

// WriteDOT writes the <IMG/> element to the writer.
func (img *IMG) WriteDOT(w io.Writer) error {
	if _, err := io.WriteString(w, "<IMG"); err != nil {
		return err
	}
	// SRC is required for IMG
	if img.Src == "" {
		return fmt.Errorf("IMG element requires SRC attribute")
	}
	writeAttrString(w, "SRC", img.Src)
	writeAttrHtmLikeImgScale(w, "SCALE", img.Scale)

	_, err := io.WriteString(w, "/>") // Self-closing
	return err
}

// Style elements (B, I, U, O, SUB, SUP, S)
type StyleElement struct {
	Tag      string
	Children []HtmLikeElement
}

func newStyleElement(tag string, children ...HtmLikeElement) *StyleElement {
	return &StyleElement{Tag: tag, Children: children}
}

// NewB creates a <B> element.
func NewB(children ...HtmLikeElement) *StyleElement { return newStyleElement("B", children...) }

// NewI creates an <I> element.
func NewI(children ...HtmLikeElement) *StyleElement { return newStyleElement("I", children...) }

// NewU creates a <U> element.
func NewU(children ...HtmLikeElement) *StyleElement { return newStyleElement("U", children...) }

// NewO creates an <O> element.
func NewO(children ...HtmLikeElement) *StyleElement { return newStyleElement("O", children...) }

// NewSUB creates a <SUB> element.
func NewSUB(children ...HtmLikeElement) *StyleElement { return newStyleElement("SUB", children...) }

// NewSUP creates a <SUP> element.
func NewSUP(children ...HtmLikeElement) *StyleElement { return newStyleElement("SUP", children...) }

// NewS creates an <S> element.
func NewS(children ...HtmLikeElement) *StyleElement { return newStyleElement("S", children...) }

// AppendChildren adds children to the style element.
func (s *StyleElement) AppendChildren(children ...HtmLikeElement) *StyleElement {
	s.Children = append(s.Children, children...)
	return s
}

// WriteDOT writes the style element and its children to the writer.
func (s *StyleElement) WriteDOT(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "<%s>", s.Tag); err != nil {
		return err
	}
	for _, child := range s.Children {
		if err := child.WriteDOT(w); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "</%s>", s.Tag)
	return err
}

// Horizontal and Vertical Rules (HR, VR)
type RuleElement struct {
	Tag string
}

func newRuleElement(tag string) *RuleElement {
	return &RuleElement{Tag: tag}
}

// NewHR creates an <HR/> element.
func NewHR() *RuleElement { return newRuleElement("HR") }

// NewVR creates a <VR/> element.
func NewVR() *RuleElement { return newRuleElement("VR") }

// WriteDOT writes the rule element to the writer.
func (r *RuleElement) WriteDOT(w io.Writer) error {
	_, err := fmt.Fprintf(w, "<%s/>", r.Tag) // Self-closing
	return err
}

// NewHtmLike creates a new HtmLike label with the given root elements.
func NewHtmLike(elements ...HtmLikeElement) *HtmLike {
	return &HtmLike{Elements: elements}
}

// AppendElements adds root elements to the HtmLike label.
func (h *HtmLike) AppendElements(elements ...HtmLikeElement) *HtmLike {
	h.Elements = append(h.Elements, elements...)
	return h
}

// String generates the Graphviz HTML-like label string representation (<...>)
func (h *HtmLike) String() string {
	if len(h.Elements) == 0 {
		return "<>" // Empty label?
	}

	var buf bytes.Buffer
	// Write the opening angle bracket
	/*
		if _, err := buf.WriteString("<"); err != nil {
			// Should not happen with bytes.Buffer
			return fmt.Sprintf("Error writing start tag: %v", err)
		}
	*/
	// Write each root element
	for _, el := range h.Elements {
		if err := el.WriteDOT(&buf); err != nil {
			// Should not happen with bytes.Buffer or simple writes
			return fmt.Sprintf("Error writing element: %v", err)
		}
	}

	/*
		// Write the closing angle bracket
		if _, err := buf.WriteString(">"); err != nil {
			// Should not happen with bytes.Buffer
			return fmt.Sprintf("Error writing end tag: %v", err)
		}
	*/
	return buf.String()
}

// Example usage:
/*
import (
	"fmt"
	"github.com/your_fork/dot/htmllike" // Adjust import path
)

func ExampleHtmLike() {
	// Create a simple table with two rows and two cells each
	label := htmllike.NewHtmLike(
		htmllike.NewTable().
			SetBorder(1).
			SetCellPadding(5).
			SetCellSpacing(0). // Example of issue mentioned in spec
			AppendChildren(
				htmllike.NewTR().AppendChildren(
					htmllike.NewTD().SetBGColor("lightblue").AppendChildren(htmllike.NewText("Header 1")),
					htmllike.NewTD().SetBGColor("lightgreen").AppendChildren(htmllike.NewText("Header 2")),
				),
				htmllike.NewHR(), // Horizontal rule between rows
				htmllike.NewTR().AppendChildren(
					htmllike.NewTD().SetPort("cell_a").AppendChildren(
						htmllike.NewText("Row 1, Cell A"),
						htmllike.NewBR(), // Line break within a cell
						htmllike.NewFONT().SetColor("red").AppendChildren(htmllike.NewText("Styled text")),
					),
					htmllike.NewVR(), // Vertical rule between cells (less common in TR)
					htmllike.NewTD().SetRowspan(2).AppendChildren(htmllike.NewText("Row 1, Cell B (Spans 2 rows)")),
				),
				htmllike.NewTR().AppendChildren(
					htmllike.NewTD().SetColspan(1).AppendChildren(htmllike.NewText("Row 2, Cell A")),
				),
			),
	)

	// This string can now be assigned to a label attribute in your dot graph.
	dotLabelString := label.String()
	fmt.Println(dotLabelString)

	// Example with just styled text
	styledLabel := htmllike.NewHtmLike(
		htmllike.NewB(htmllike.NewText("Bold ")),
		htmllike.NewI(htmllike.NewText("Italic")),
	)
	fmt.Println(styledLabel.String())

	// Example with an image
	imageLabel := htmllike.NewHtmLike(
		htmllike.NewTable().AppendChildren(
			htmllike.NewTR().AppendChildren(
				htmllike.NewTD().AppendChildren(htmllike.NewIMG("path/to/image.png").SetScale(htmllike.ScaleBOTH)),
				htmllike.NewTD().AppendChildren(htmllike.NewText("Image Description")),
			),
		),
	)
	fmt.Println(imageLabel.String())
}
*/
