package helpers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/juju/errors"
)

type TestMarshable struct {
	data        []byte
	marshalErr  error
	validateErr error
}

func (m *TestMarshable) MarshalBinary() ([]byte, error)         { return m.data, m.marshalErr }
func (m *TestMarshable) Validate(formats strfmt.Registry) error { return m.validateErr }

func TestWriteRes(t *testing.T) {
	t.Run("writes ISE if validation fails", func(t *testing.T) {
		rr := httptest.NewRecorder()
		m := &TestMarshable{
			validateErr: errors.New("validation error"),
		}

		WriteRes(rr, m)

		if got, want := rr.Code, http.StatusInternalServerError; got != want {
			t.Errorf("wrong status code. got: %d, want: %d", got, want)
		}
	})

	t.Run("writes ISE if cannot marshall", func(t *testing.T) {
		rr := httptest.NewRecorder()
		wantErr := errors.New("marshall error")
		m := &TestMarshable{
			marshalErr: wantErr,
		}

		WriteRes(rr, m)

		if got, want := rr.Code, http.StatusInternalServerError; got != want {
			t.Errorf("wrong status code. got: %d, want: %d", got, want)
		}
	})

	t.Run("panics if cannot write response", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("this was expected to panic, but it didn't")
			}
		}()

		var rr *httptest.ResponseRecorder
		m := &TestMarshable{}

		WriteRes(rr, m)
	})

	t.Run("writes response", func(t *testing.T) {})
}
