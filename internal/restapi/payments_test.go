package restapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/getlantern/deepcopy"
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/j-fuentes/payments/internal/fixtures"
	"github.com/j-fuentes/payments/internal/store"
	"github.com/j-fuentes/payments/pkg/models"
)

func copyPayment(p *models.Payment) *models.Payment {
	var result models.Payment
	deepcopy.Copy(&result, p)
	result.ID = strfmt.UUID(uuid.New().String())
	return &result
}

func TestGetPayments(t *testing.T) {
	file := "single.json"
	fixture, err := fixtures.LoadPayments(file)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if l := len(fixture.Data); l != 1 {
		t.Fatalf("expected to load just one payment from %s, but found %d", file, l)
	}

	canonicalProject := fixture.Data[0]
	canonicalProject.Attributes.Amount = "100.0"

	p1 := copyPayment(canonicalProject)
	p2 := copyPayment(canonicalProject)
	p3 := copyPayment(canonicalProject)
	p4 := copyPayment(canonicalProject)
	p5 := copyPayment(canonicalProject)
	p6 := copyPayment(canonicalProject)

	payments := []*models.Payment{p1, p2, p3, p4, p5, p6}

	t.Run("returns a valid list of payments", func(t *testing.T) {
		sv := NewPaymentsServer(store.NewVolatilePaymentsStore(payments), "http://localhost:3000")

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

		if got, want := res.Data, payments; !reflect.DeepEqual(got, want) {
			t.Errorf("Data does not match. got: %+v, want: %+v", got, want)
		}
	})

	t.Run("returns a valid list of payments filtering by organisation-id", func(t *testing.T) {
		orgID := strfmt.UUID(uuid.New().String())
		payments[0].OrganisationID = orgID
		payments[2].OrganisationID = orgID

		sv := NewPaymentsServer(store.NewVolatilePaymentsStore(payments), "http://localhost:3000")

		req, err := http.NewRequest("GET", fmt.Sprintf("/payments?organisation-id=%s", orgID), nil)
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

		if got, want := res.Data, []*models.Payment{payments[0], payments[2]}; !reflect.DeepEqual(got, want) {
			t.Errorf("Data does not match. got: %+v, want: %+v", got, want)
		}
	})

	t.Run("returns a valid list of payments filtering by min-amount", func(t *testing.T) {
		payments[0].Attributes.Amount = "500.3"
		payments[2].Attributes.Amount = "100.1"

		sv := NewPaymentsServer(store.NewVolatilePaymentsStore(payments), "http://localhost:3000")

		req, err := http.NewRequest("GET", fmt.Sprintf("/payments?min-amount=%s", "100.1"), nil)
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

		if got, want := res.Data, []*models.Payment{payments[0], payments[2]}; !reflect.DeepEqual(got, want) {
			t.Errorf("Data does not match. got: %+v, want: %+v", got, want)
		}
	})

	t.Run("returns a valid list of payments filtering by max-amount", func(t *testing.T) {
		payments[0].Attributes.Amount = "50.3"
		payments[2].Attributes.Amount = "10.1"

		sv := NewPaymentsServer(store.NewVolatilePaymentsStore(payments), "http://localhost:3000")

		req, err := http.NewRequest("GET", fmt.Sprintf("/payments?max-amount=%s", "50.3"), nil)
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

		if got, want := res.Data, []*models.Payment{payments[0], payments[2]}; !reflect.DeepEqual(got, want) {
			t.Errorf("Data does not match. got: %+v, want: %+v", got, want)
		}
	})
}
