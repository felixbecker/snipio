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
	ErrFileImport error = fmt.Errorf("error importing the draw io file")
	// ErrFileParsing will be returned when the draw io file could not be parsed
	ErrFileParsing error = fmt.Errorf("error parsing the draw io xml")

	// ErrNoLayersFound will be returned when no layer was found in the model
	ErrNoLayersFound error = fmt.Errorf("error no layers found")

	// ErrLayerNotFound will be returned when the given layer is not found
	ErrLayerNotFound error = fmt.Errorf("error layer not found")

	// ErrExportingXML will be returned when marshalling the struct into xml does not work
	ErrExportingXML error = fmt.Errorf("error generating the xml for the export")

	// ErrNoValidLayerName will be returnd if the layer name is not valid
	ErrNoValidLayerName error = fmt.Errorf("error no valid layer name")

	// ErrNoCells will be returned if the cells are nil or empty
	ErrNoCells error = fmt.Errorf("error cells are nil and should have a value")

	// ErrNoID will be returned if the id is a empty string
	ErrNoID error = fmt.Errorf("error id is empty and should have a value")
)

//go:embed draft.xml
var classificationString []byte

// App holds all data for an App
type App struct {
	layers []LayerInfo
	model  model
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

func importDrawing(filename string) ([]LayerInfo, *model, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, ErrFileImport
	}

	if strings.Contains(string(data), "mxfile") {
		data, err = importMxFile(data)
		if err != nil {
			return nil, nil, ErrFileImport
		}
	}

	var m model
	err = xml.Unmarshal(data, &m)
	if err != nil {
		fmt.Println(err)
		return nil, nil, ErrFileParsing
	}

	layers := []LayerInfo{}
	for i, c := range m.Cells {
		for _, attr := range c.Attributes {
			if strings.ToLower(attr.Name.Local) == "id" {
				m.Cells[i].ID = attr.Value
			}
			if strings.ToLower(attr.Name.Local) == "parent" {

				m.Cells[i].Parent = attr.Value
			}
			if strings.ToLower(attr.Name.Local) == "value" {
				m.Cells[i].Value = attr.Value
			}

			if m.Cells[i].Parent == "0" {

				if m.Cells[i].Value == "" {
					m.Cells[i].Value = "Background"
				}
				layers = append(layers, LayerInfo{
					ID:   m.Cells[i].ID,
					Name: m.Cells[i].Value,
					Idx:  i,
				})
			}

		}
	}

	return layers, &m, nil

}

// ImportDrawing imports the draw io drawing
func (a *App) ImportDrawing(filename string) error {

	layers, m, err := importDrawing(filename)
	if err != nil {
		return err
	}
	a.layers = layers
	a.model = *m
	return nil
}

// RemoveLayerByName removes a layer by a given name
func (a *App) RemoveLayerByName(name string, outputFilename string) error {

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

	err = writeFile(outputFilename, a.model)
	if err != nil {
		return err
	}
	return nil

}

// ExtractLayerByName exports only the layer by given name
func (a *App) ExtractLayerByName(name string, outputFile string) error {

	id := a.layerID(name)
	if id == "" {
		return ErrLayerNotFound
	}
	cells, err := keepElementsWithID(a.model.Cells, id)
	if err != nil {
		return err
	}
	a.model.Cells = cells
	err = writeFile(outputFile, a.model)
	if err != nil {
		return err
	}
	return nil
}

// Merge takes all layers and merges them onto the imported files
func (a *App) Merge(filenameToBeMerged string, outputFilename string) error {

	if len(filenameToBeMerged) == 0 {
		return fmt.Errorf("no file to be merged found")
	}
	_, m, err := importDrawing(filenameToBeMerged)
	if err != nil {
		return err
	}

	m.Cells = findAndDelete(m.Cells, "0")
	m.Cells = findAndDelete(m.Cells, "1")
	a.model.Cells = append(a.model.Cells, m.Cells...)
	err = writeFile(outputFilename, a.model)
	if err != nil {
		return err
	}
	return nil
}

//Classify marks a document as draft
func (a *App) Classify(targetfilename string) error {

	classification, err := makeClassificationLabel(string(classificationString))
	if err != nil {
		return err
	}
	a.model.Cells = append(a.model.Cells, classification...)
	err = writeFile(targetfilename, a.model)
	if err != nil {
		return err
	}
	return err
}

func (a *App) layerID(name string) string {
	for _, v := range a.layers {
		if name == v.Name {
			return v.ID
		}
	}
	return "" //not found.
}

func checkCellsAndIDcontainsValues(s []cell, id string) error {
	if s == nil {
		return ErrNoCells
	}
	if len(s) == 0 {
		return ErrNoCells
	}
	if len(id) == 0 {

		return ErrNoID
	}
	return nil
}

func keepElementsWithID(s []cell, id string) ([]cell, error) {

	err := checkCellsAndIDcontainsValues(s, id)
	if err != nil {
		return nil, err
	}

	cells := []cell{{XMLName: xml.Name{Local: "mxCell"}, Attributes: []xml.Attr{
		{Name: xml.Name{Local: "id"}, Value: "0"}}}}

	for _, c := range s {

		if c.ID == id || c.Parent == id {
			cells = append(cells, c)
		}
	}

	return cells, nil
}

func removeElementsWithID(s []cell, id string) ([]cell, error) {

	err := checkCellsAndIDcontainsValues(s, id)
	if err != nil {
		return nil, err
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
	if err != nil {
		return err
	}

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

func findAndDelete(s []cell, id string) []cell {

	index := 0
	for _, i := range s {
		if i.ID != id {
			s[index] = i
			index++
		}
	}
	return s[:index]
}
