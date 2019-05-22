package restapi

import (
	"net/http"
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"
	"bytes"

	"github.com/getlantern/deepcopy"
	"github.com/go-openapi/strfmt"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"github.com/j-fuentes/payments/internal/fixtures"
	"github.com/j-fuentes/payments/internal/store"
	"github.com/j-fuentes/payments/pkg/models"
)

func paymentFromFixture(t *testing.T) *models.Payment {
	file := "single.json"
	fixture, err := fixtures.LoadPayments(file)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if l := len(fixture.Data); l != 1 {
		t.Fatalf("expected to load just one payment from %s, but found %d", file, l)
	}

	return fixture.Data[0]
}

func newPaymentFrom(p *models.Payment) *models.Payment {
	var result models.Payment
	deepcopy.Copy(&result, p)
	result.ID = strfmt.UUID(uuid.New().String())
	return &result
}

func TestGetPayments(t *testing.T) {
	canonicalProject := paymentFromFixture(t)
	canonicalProject.Attributes.Amount = "100.0"

	p1 := newPaymentFrom(canonicalProject)
	p2 := newPaymentFrom(canonicalProject)
	p3 := newPaymentFrom(canonicalProject)
	p4 := newPaymentFrom(canonicalProject)
	p5 := newPaymentFrom(canonicalProject)
	p6 := newPaymentFrom(canonicalProject)

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

func TestGetPayment(t *testing.T) {
	canonicalProject := paymentFromFixture(t)

	p1 := newPaymentFrom(canonicalProject)
	p2 := newPaymentFrom(canonicalProject)

	payments := []*models.Payment{p1, p2}
	fmt.Println(payments)

	t.Run("finds and returns a payment", func(t *testing.T) {
		sv := NewPaymentsServer(store.NewVolatilePaymentsStore(payments), "http://localhost:3000")

		req, err := http.NewRequest("GET", fmt.Sprintf("/payment/%s", p2.ID), nil)
		if err != nil {
			t.Fatal(err)
		}

		req = mux.SetURLVars(req, map[string]string{
			"id": p2.ID.String(),
		})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sv.GetPayment)

		handler.ServeHTTP(rr, req)

		if got, want := rr.Code, http.StatusOK; got != want {
			t.Errorf("handler returned wrong status code. got: %d, want: %d", got, want)
		}

		var res models.Payment
		err = res.UnmarshalBinary(rr.Body.Bytes())
		if err != nil {
			t.Errorf("Cannot unmarshal response: %+v", err)
		}
		err = res.Validate(nil)
		if err != nil {
			t.Errorf("Validation failed for response: %+v", err)
		}

		if got, want := &res, p2; !reflect.DeepEqual(got, want) {
			t.Errorf("Data does not match. got: %+v, want: %+v", got, want)
		}
	})

	t.Run("returns 404 if not found", func(t *testing.T) {
		sv := NewPaymentsServer(store.NewVolatilePaymentsStore(payments), "http://localhost:3000")

		uuid := strfmt.UUID("deadbeef")
		req, err := http.NewRequest("GET", fmt.Sprintf("/payment/%s", uuid), nil)
		if err != nil {
			t.Fatal(err)
		}

		req = mux.SetURLVars(req, map[string]string{
			"id": uuid.String(),
		})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sv.GetPayment)

		handler.ServeHTTP(rr, req)

		if got, want := rr.Code, http.StatusNotFound; got != want {
			t.Errorf("handler returned wrong status code. got: %d, want: %d", got, want)
		}
	})
}

func TestDeletePayment(t *testing.T) {
	canonicalProject := paymentFromFixture(t)

	p1 := newPaymentFrom(canonicalProject)
	p2 := newPaymentFrom(canonicalProject)

	payments := []*models.Payment{p1, p2}
	fmt.Println(payments)

	t.Run("deletes a payment", func(t *testing.T) {
		sv := NewPaymentsServer(store.NewVolatilePaymentsStore(payments), "http://localhost:3000")

		req, err := http.NewRequest("DELETE", fmt.Sprintf("/payment/%s", p2.ID), nil)
		if err != nil {
			t.Fatal(err)
		}

		req = mux.SetURLVars(req, map[string]string{
			"id": p2.ID.String(),
		})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sv.DeletePayment)

		handler.ServeHTTP(rr, req)

		if got, want := rr.Code, http.StatusOK; got != want {
			t.Errorf("handler returned wrong status code. got: %d, want: %d", got, want)
		}

		var res models.Empty
		err = res.UnmarshalBinary(rr.Body.Bytes())
		if err != nil {
			t.Errorf("Cannot unmarshal response: %+v", err)
		}
		err = res.Validate(nil)
		if err != nil {
			t.Errorf("Validation failed for response: %+v", err)
		}
	})

	t.Run("returns 404 if not found", func(t *testing.T) {
		sv := NewPaymentsServer(store.NewVolatilePaymentsStore(payments), "http://localhost:3000")

		uuid := strfmt.UUID("deadbeef")
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/payment/%s", uuid), nil)
		if err != nil {
			t.Fatal(err)
		}

		req = mux.SetURLVars(req, map[string]string{
			"id": uuid.String(),
		})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sv.DeletePayment)

		handler.ServeHTTP(rr, req)

		if got, want := rr.Code, http.StatusNotFound; got != want {
			t.Errorf("handler returned wrong status code. got: %d, want: %d", got, want)
		}
	})
}

func TestUpdatePayment(t *testing.T) {
	canonicalProject := paymentFromFixture(t)

	p1 := newPaymentFrom(canonicalProject)
	p2 := newPaymentFrom(canonicalProject)

	payments := []*models.Payment{p1, p2}
	fmt.Println(payments)

	t.Run("updates a payment", func(t *testing.T) {
		sv := NewPaymentsServer(store.NewVolatilePaymentsStore(payments), "http://localhost:3000")

		bb, err := p2.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}


		req, err := http.NewRequest("UPDATE", fmt.Sprintf("/payment/%s", p2.ID), bytes.NewReader(bb))
		if err != nil {
			t.Fatal(err)
		}

		req = mux.SetURLVars(req, map[string]string{
			"id": p2.ID.String(),
		})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sv.UpdatePayment)

		handler.ServeHTTP(rr, req)

		if got, want := rr.Code, http.StatusOK; got != want {
			t.Errorf("handler returned wrong status code. got: %d, want: %d", got, want)
		}

		var res models.Empty
		err = res.UnmarshalBinary(rr.Body.Bytes())
		if err != nil {
			t.Errorf("Cannot unmarshal response: %+v", err)
		}
		err = res.Validate(nil)
		if err != nil {
			t.Errorf("Validation failed for response: %+v", err)
		}
	})

	t.Run("returns 404 if not found", func(t *testing.T) {
		sv := NewPaymentsServer(store.NewVolatilePaymentsStore(payments), "http://localhost:3000")

		p := newPaymentFrom(p1)
		p.ID = strfmt.UUID(uuid.New().String())
		bb, err := p.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}


		req, err := http.NewRequest("UPDATE", fmt.Sprintf("/payment/%s", p.ID), bytes.NewReader(bb))
		if err != nil {
			t.Fatal(err)
		}

		req = mux.SetURLVars(req, map[string]string{
			"id": p.ID.String(),
		})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sv.UpdatePayment)

		handler.ServeHTTP(rr, req)

		if got, want := rr.Code, http.StatusNotFound; got != want {
			t.Errorf("handler returned wrong status code. got: %d, want: %d", got, want)
		}
	})

	t.Run("returns 400 if cannot unmarshal payment", func(t *testing.T) {
		sv := NewPaymentsServer(store.NewVolatilePaymentsStore(payments), "http://localhost:3000")

		p := newPaymentFrom(p1)

		req, err := http.NewRequest("UPDATE", fmt.Sprintf("/payment/%s", p.ID), bytes.NewReader([]byte{}))
		if err != nil {
			t.Fatal(err)
		}

		req = mux.SetURLVars(req, map[string]string{
			"id": p.ID.String(),
		})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sv.UpdatePayment)

		handler.ServeHTTP(rr, req)

		if got, want := rr.Code, http.StatusBadRequest; got != want {
			t.Errorf("handler returned wrong status code. got: %d, want: %d", got, want)
		}
	})

	t.Run("returns 400 if cannot validate payment", func(t *testing.T) {
		sv := NewPaymentsServer(store.NewVolatilePaymentsStore(payments), "http://localhost:3000")

		p := newPaymentFrom(p1)
		p.Type = ""
		bb, err := p.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("UPDATE", fmt.Sprintf("/payment/%s", p.ID), bytes.NewReader(bb))
		if err != nil {
			t.Fatal(err)
		}

		req = mux.SetURLVars(req, map[string]string{
			"id": p.ID.String(),
		})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(sv.UpdatePayment)

		handler.ServeHTTP(rr, req)

		if got, want := rr.Code, http.StatusBadRequest; got != want {
			t.Errorf("handler returned wrong status code. got: %d, want: %d", got, want)
		}
	})
}
