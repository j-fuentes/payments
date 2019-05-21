package fixtures

import (
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/j-fuentes/payments/pkg/models"
)

func testdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "testdata")
}

// LoadPayments reads a Payments struct from a JSON file.
func LoadPayments(name string) (*models.Payments, error) {
	path := filepath.Join(testdataDir(), name)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var payments models.Payments
	err = payments.UnmarshalBinary(bytes)
	if err != nil {
		return nil, err
	}

	return &payments, nil
}
