package app

import (
	"encoding/xml"
	"errors"
	"testing"

	"github.com/matryer/is"
)

func Test_keepElementsWithID_fails_if_nil_and_empty_string_is_passed_in(t *testing.T) {

	is := is.New(t)

	cells, err := keepElementsWithID(nil, "")

	is.True(errors.Is(err, ErrNoCells))
	is.True(cells == nil)
}

func Test_keepElementsWithID_fails_if_empty_cell_slice_is_passed_in(t *testing.T) {

	is := is.New(t)

	cells, err := keepElementsWithID([]cell{}, "")

	is.True(errors.Is(err, ErrNoCells))
	is.True(cells == nil)
}

func Test_keepElementsWithID_fails_with_empty_ID(t *testing.T) {
	is := is.New(t)
	cs := []cell{{XMLName: xml.Name{Local: "mxCell"}, Attributes: []xml.Attr{
		{Name: xml.Name{Local: "id"}, Value: "0"}}}}
	cells, err := keepElementsWithID(cs, "")
	is.True(errors.Is(err, ErrNoID))
	is.True(cells == nil)

}
