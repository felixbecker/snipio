package app

import (
	"bytes"
	"compress/flate"

	// embed should embed
	_ "embed"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
)

var (
	// ErrFileImport will be returned when the draw io file could not be imported
	ErrFileImport error = fmt.Errorf("Error importing the draw io file")
	// ErrFileParsing will be returned when the draw io file could not be parsed
	ErrFileParsing error = fmt.Errorf("Error parsing the draw io xml")

	// ErrNoLayersFound will be returned when no layer was found in the model
	ErrNoLayersFound error = fmt.Errorf("Error no layers found")

	// ErrLayerNotFound will be returned when the given layer is not found
	ErrLayerNotFound error = fmt.Errorf("Error layer not found")

	// ErrExportingXML will be returned when marshalling the struct into xml does not work
	ErrExportingXML error = fmt.Errorf("Error generating the xml for the export")

	// ErrNoValidLayerName will be returnde if the layer name is not valid
	ErrNoValidLayerName error = fmt.Errorf("Error no valid layer name")
)

//go:embed draft.xml
var classificationString []byte

// App holds all data for an App
type App struct {
	classification []cell
	layers         []LayerInfo
	model          model
}

// New creates a new App
func New() *App {

	return &App{}
}

func makeClassificationLabel(templateString string) ([]cell, error) {
	type root struct {
		Cells []cell `xml:"mxCell"`
	}
	classification := root{}
	err := xml.Unmarshal(classificationString, &classification)
	if err != nil {
		return nil, err
	}
	return classification.Cells, err

}

// Layers returns a list of layer names with its ids
func (a *App) Layers() ([]LayerInfo, error) {

	if len(a.layers) == 0 {
		return nil, ErrNoLayersFound
	}
	return a.layers, nil

}

func writeFile(filename string, model model) error {

	bts, err := xml.Marshal(model)
	if err != nil {
		fmt.Println(err)
		return ErrExportingXML
	}

	err = ioutil.WriteFile(filename, bts, 0644)
	if err != nil {
		return err
	}
	fmt.Println("saved data to file: ", filename)
	return nil

}

// ImportDrawing imports the draw io drawing
func (a *App) ImportDrawing(filename string) error {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return ErrFileImport
	}

	if strings.Contains(string(data), "mxfile") {
		data, err = importMxFile(data)
		if err != nil {
			return ErrFileImport
		}
	}

	err = xml.Unmarshal(data, &a.model)
	if err != nil {
		fmt.Println(err)
		return ErrFileParsing
	}

	layers := []LayerInfo{}
	for i, c := range a.model.Cells {
		for _, attr := range c.Attributes {
			if strings.ToLower(attr.Name.Local) == "id" {
				a.model.Cells[i].ID = attr.Value
			}
			if strings.ToLower(attr.Name.Local) == "parent" {

				a.model.Cells[i].Parent = attr.Value
			}
			if strings.ToLower(attr.Name.Local) == "value" {
				a.model.Cells[i].Value = attr.Value
			}

			if a.model.Cells[i].Parent == "0" {

				if a.model.Cells[i].Value == "" {
					a.model.Cells[i].Value = "Background"
				}
				layers = append(layers, LayerInfo{
					ID:   a.model.Cells[i].ID,
					Name: a.model.Cells[i].Value,
					Idx:  i,
				})
			}

		}
	}

	a.layers = layers

	a.classification, err = makeClassificationLabel(string(classificationString))
	if err != nil {
		fmt.Println(err)
		return ErrFileParsing
	}
	return nil
}

// RemoveLayerByName removes a layer by a given name
func (a *App) RemoveLayerByName(name string, targetfilename string) error {

	if len(name) == 0 {
		return ErrNoValidLayerName
	}

	id := a.layerID(name)
	if len(id) == 0 {
		return ErrLayerNotFound
	}

	cells, err := removeElementsWithID(a.model.Cells, id)

	if err != nil {
		return err
	}
	a.model.Cells = cells

	err = writeFile(targetfilename, a.model)
	if err != nil {
		return err
	}
	return nil

}

// ExtractLayerByName exports only the layer by given name
func (a *App) ExtractLayerByName(name string, targetFilename string) error {

	id := a.layerID(name)
	if id == "" {
		return ErrLayerNotFound
	}
	a.model.Cells = keepElementsWithID(a.model.Cells, id)
	err := writeFile(targetFilename, a.model)
	if err != nil {
		return err
	}
	return nil
}

//Classify marks a document as draft
func (a *App) Classify(targetfilename string) {

	a.model.Cells = append(a.model.Cells, a.classification...)
	writeFile(targetfilename, a.model)
}

func (a *App) layerID(name string) string {
	for _, v := range a.layers {
		if name == v.Name {
			return v.ID
		}
	}
	return "" //not found.
}

func keepElementsWithID(s []cell, id string) []cell {

	cells := []cell{{XMLName: xml.Name{Local: "mxCell"}, Attributes: []xml.Attr{
		{Name: xml.Name{Local: "id"}, Value: "0"}}}}

	for _, c := range s {

		if c.ID == id || c.Parent == id {
			cells = append(cells, c)
		}
	}

	return cells
}
func removeElementsWithID(s []cell, id string) ([]cell, error) {

	if len(id) == 0 {
		return nil, fmt.Errorf("no id")
	}

	cells := []cell{{XMLName: xml.Name{Local: "mxCell"}, Attributes: []xml.Attr{
		{Name: xml.Name{Local: "id"}, Value: "0"}}}}
	for _, c := range s {

		if c.ID != id && c.Parent != id {
			cells = append(cells, c)
		}
	}

	return cells, nil
}

// UnpackFile unpacks an mxfile and etracts the xml
func (a *App) UnpackFile(filename string, outputFile string) error {

	bts, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if !strings.Contains(string(bts), "mxfile") {
		return fmt.Errorf("file is not a mxfile")
	}

	bts, err = importMxFile(bts)
	if err != nil {
		return err
	}

	if len(outputFile) == 0 {
		fmt.Println(string(bts))
		return nil
	}

	err = ioutil.WriteFile(outputFile, bts, 0644)

	return nil
}

func importMxFile(data []byte) ([]byte, error) {

	fileData := mxFile{}
	err := xml.Unmarshal(data, &fileData)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("YIKES")
	}

	bts, err := base64.StdEncoding.DecodeString(fileData.Diagram)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("YIKES")
	}

	r := flate.NewReader(bytes.NewReader(bts))
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("YIKES")
	}
	enflated, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("YIKES")
	}
	decodedValue, err := url.QueryUnescape(string(enflated))
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("YIKES")
	}

	return []byte(decodedValue), nil

}
