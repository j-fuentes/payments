package helpers

import (
	"net/http"

	"github.com/go-openapi/strfmt"
	"github.com/golang/glog"
	"github.com/j-fuentes/payments/pkg/models"
)

// Marshable is something that can be rendered into JSON according to the schema.
type Marshable interface {
	MarshalBinary() ([]byte, error)
	Validate(formats strfmt.Registry) error
}

// WriteRes writes in an http response a Marshable
func WriteRes(w http.ResponseWriter, m Marshable) {
	err := m.Validate(nil)
	if err != nil {
		glog.Errorf("Validation error: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bb, err := m.MarshalBinary()
	if err != nil {
		glog.Errorf("Cannot marshall object: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bb)
	if err != nil {
		panic(err)
	}

	return
}

func WriteError(w http.ResponseWriter, code int, err error) {
	m := &models.Error{
		Code:    int64(code),
		Message: err.Error(),
	}

	w.WriteHeader(code)

	WriteRes(w, m)
}
