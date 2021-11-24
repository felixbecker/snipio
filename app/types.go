package app

import "encoding/xml"

// LayerInfo holds the information about a layer name and its id
type LayerInfo struct {
	Name string
	ID   string
	Idx  int
}

// M test
type model struct {
	XMLName         xml.Name   `xml:"mxGraphModel"`
	ModelAttributes []xml.Attr `xml:",any,attr"`
	Cells           []cell     `xml:"root>mxCell"`
}

// C test
type cell struct {
	XMLName    xml.Name   `xml:"mxCell"`
	Attributes []xml.Attr `xml:",any,attr"`
	Content    string     `xml:",innerxml"`
	ID         string     `xml:"-"`
	Parent     string     `xml:"-"`
	Value      string     `xml:"-"`
}

type mxFile struct {
	Diagram string `xml:"diagram"`
}
