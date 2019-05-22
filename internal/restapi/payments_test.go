package restapi

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/j-fuentes/payments/internal/fixtures"
	"github.com/j-fuentes/payments/internal/store"
	"github.com/j-fuentes/payments/pkg/models"
)

func TestGetPayments(t *testing.T) {
	payments, err := fixtures.LoadPayments("payments.json")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	sv := NewPaymentsServer(store.NewVolatilePaymentsStore(payments.Data), "http://localhost:3000")

	t.Run("returns a valid list of payments", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/payments", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sv.GetPayments)

		handler.ServeHTTP(rr, req)

		if got, want := rr.Code, http.StatusOK; got != want {
			t.Errorf("handler returned wrong status code. got: %d, want: %d", got, want)
		}

		var res models.Payments
		err = res.UnmarshalBinary(rr.Body.Bytes())
		if err != nil {
			t.Errorf("Cannot unmarshal response: %+v", err)
		}
		err = res.Validate(nil)
		if err != nil {
			t.Errorf("Validation failed for response: %+v", err)
		}

		if got, want := res.Data, payments.Data; !reflect.DeepEqual(got, want) {
			t.Errorf("Data does not match. got: %+v, want: %+v", got, want)
		}

		// TODO: Validate Link once we define what it should be
	})
}
