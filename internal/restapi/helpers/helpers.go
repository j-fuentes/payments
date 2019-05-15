package helpers

import "net/http"

// Marshable is something that can be rendered into JSON accoring to the schema.
type Marshable interface {
	MarshalBinary() ([]byte, error)
}

// WriteRes writes in an http response a Marshable
func WriteRes(w http.ResponseWriter, m Marshable) error {
	bb, err := m.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = w.Write(bb)
	if err != nil {
		return err
	}

	return nil
}
