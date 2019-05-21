package helpers

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/golang/glog"
	"github.com/juju/errors"
)

// Marshable is something that can be rendered into JSON according to the schema.
type Marshable interface {
	MarshalBinary() ([]byte, error)
	Validate(formats strfmt.Registry) error
}

// WriteRes writes in an http response a Marshable
func WriteRes(w http.ResponseWriter, m Marshable) error {
	err := m.Validate(nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return errors.BadRequestf("%+v", err)
	}

	bb, err := m.MarshalBinary()
	if err != nil {
		glog.Errorf("Cannot marshall object: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	_, err = w.Write(bb)
	if err != nil {
		panic(err)
	}

	return nil
}

// TODO: define what this should return
func GenerateLink(path string) string {
	return fmt.Sprintf("https://domain%s", path)
}
