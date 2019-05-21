package fixtures

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadPayments asserts both the function LoadPayments works as expected and also that the data in the testdata dir is valid according to the models.
func TestLoadPayments(t *testing.T) {
	err := filepath.Walk(testdataDir(), func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".json" {
			return nil
		}

		payments, err := LoadPayments(filepath.Base(path))
		if err != nil {
			t.Fatalf("unexpected error loading %q: %+v", path, err)
		}

		err = payments.Validate(nil)
		if err != nil {
			t.Errorf("validation failed on root Payments object in %q: %+v", path, err)

			// in case of validation error in main object, iterate to get the index if the payment in the array, useful to debug
			for i, p := range payments.Data {
				err = p.Validate(nil)
				if err != nil {
					t.Errorf("validation failed for payment #%d in %q: %+v", i, path, err)
				}
			}

		}

		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}
}
