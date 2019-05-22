package store_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/getlantern/deepcopy"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/j-fuentes/payments/internal/fixtures"
	"github.com/j-fuentes/payments/internal/store"
	"github.com/j-fuentes/payments/pkg/models"
	"github.com/juju/errors"
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

func clonePayment(p *models.Payment) *models.Payment {
	var result models.Payment
	deepcopy.Copy(&result, p)
	return &result
}

func newPaymentFrom(p *models.Payment) *models.Payment {
	result := clonePayment(p)
	result.ID = strfmt.UUID(uuid.New().String())
	return result
}

func TestGetPayments(t *testing.T) {
	orgID := strfmt.UUID(uuid.New().String())
	minAmount := 100.2
	maxAmount := 300.0

	canonicalProject := paymentFromFixture(t)

	p1 := newPaymentFrom(canonicalProject)

	p2 := newPaymentFrom(canonicalProject)
	p2.OrganisationID = orgID

	p3 := newPaymentFrom(canonicalProject)
	p3.Attributes.Amount = fmt.Sprintf("%.2f", minAmount+1.0)

	p4 := newPaymentFrom(canonicalProject)
	p4.OrganisationID = orgID
	p4.Attributes.Amount = fmt.Sprintf("%.2f", minAmount+1.0)

	p5 := newPaymentFrom(canonicalProject)
	p5.OrganisationID = orgID
	p5.Attributes.Amount = fmt.Sprintf("%.2f", minAmount-1.0)

	p6 := newPaymentFrom(canonicalProject)
	p5.ID = strfmt.UUID(uuid.New().String())
	p6.OrganisationID = orgID
	p6.Attributes.Amount = fmt.Sprintf("%.2f", maxAmount+1.0)

	testCases := []struct {
		name         string
		filter       store.Filter
		wantPayments []*models.Payment
	}{
		{
			"returns all the payments if empty filter",
			store.NewFilter(),
			[]*models.Payment{p1, p2, p3, p4, p5, p6},
		},
		{
			"returns all the payments of one organisation",
			store.Filter{
				OrganisationID: &orgID,
				MaxAmount:      math.Inf(1),
				MinAmount:      0,
			},
			[]*models.Payment{p2, p4, p5, p6},
		},
		{
			"returns all the payments between MinAmount and MaxAmount",
			store.Filter{
				OrganisationID: nil,
				MinAmount:      minAmount + 2.0,
				MaxAmount:      maxAmount + 2.0,
			},
			[]*models.Payment{p6},
		},
	}

	s := store.NewVolatilePaymentsStore(
		[]*models.Payment{p1, p2, p3, p4, p5, p6},
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pp, err := s.GetPayments(tc.filter)
			if err != nil {
				t.Fatalf("Got unexpected error: %+v", err)
			}

			if got, want := pp, tc.wantPayments; !reflect.DeepEqual(got, want) {
				t.Errorf("got %+v, want %+v", got, want)
			}
		})
	}
}

func TestCreatePayment(t *testing.T) {
	canonicalProject := paymentFromFixture(t)

	p1 := newPaymentFrom(canonicalProject)
	p2 := newPaymentFrom(canonicalProject)

	payments := []*models.Payment{p1, p2}
	s := store.NewVolatilePaymentsStore(payments)

	t.Run("returns BadRequest if validation for Payment fails", func(t *testing.T) {
		p := newPaymentFrom(p1)
		// Type is required to have length > 0
		p.Type = ""

		_, err := s.CreatePayment(p)
		if !errors.IsBadRequest(err) {
			t.Errorf("expected BadRequest, got %+v", err)
		}

		newPayments, err := s.GetPayments(store.NewFilter())
		if err != nil {
			t.Fatalf("%+v", err)
		}

		if got, want := newPayments, payments; !reflect.DeepEqual(got, want) {
			t.Errorf("got: %+v, want: %+v", got, want)
		}
	})

	t.Run("creates a payment", func(t *testing.T) {
		newP := newPaymentFrom(p1)

		res, err := s.CreatePayment(newP)
		if err != nil {
			t.Errorf("expected no error, got %+v", err)
		}

		newPayments, err := s.GetPayments(store.NewFilter())
		if err != nil {
			t.Fatalf("%+v", err)
		}

		if got, want := res, newP; !reflect.DeepEqual(got, want) {
			t.Errorf("Return value does not match. got: %+v, want: %+v", got, want)
		}

		if got, want := newPayments, []*models.Payment{p1, p2, newP}; !reflect.DeepEqual(got, want) {
			t.Errorf("got: %+v, want: %+v", got, want)
		}
	})
}

func TestGetPayment(t *testing.T) {
	canonicalProject := paymentFromFixture(t)

	p1 := newPaymentFrom(canonicalProject)
	p2 := newPaymentFrom(canonicalProject)

	s := store.NewVolatilePaymentsStore(
		[]*models.Payment{p1, p2},
	)

	t.Run("returns NotFound if not found", func(t *testing.T) {
		_, err := s.GetPayment(strfmt.UUID(uuid.New().String()))
		if !errors.IsNotFound(err) {
			t.Errorf("expected NotFound, got %+v", err)
		}
	})

	t.Run("returns the payment", func(t *testing.T) {
		p, err := s.GetPayment(p2.ID)
		if err != nil {
			t.Errorf("expected no error, got %+v", err)
		}

		if got, want := p, p2; !reflect.DeepEqual(got, want) {
			t.Errorf("got: %+v, want: %+v", got, want)
		}
	})
}

func TestDeletePayment(t *testing.T) {
	canonicalProject := paymentFromFixture(t)

	p1 := newPaymentFrom(canonicalProject)
	p2 := newPaymentFrom(canonicalProject)
	p3 := newPaymentFrom(canonicalProject)
	p4 := newPaymentFrom(canonicalProject)

	payments := []*models.Payment{p1, p2, p3, p4}
	s := store.NewVolatilePaymentsStore(payments)

	t.Run("returns NotFound if not found", func(t *testing.T) {
		err := s.DeletePayment(strfmt.UUID(uuid.New().String()))
		if !errors.IsNotFound(err) {
			t.Errorf("expected NotFound, got %+v", err)
		}

		newPayments, err := s.GetPayments(store.NewFilter())
		if err != nil {
			t.Fatalf("%+v", err)
		}

		if got, want := newPayments, payments; !reflect.DeepEqual(got, want) {
			t.Errorf("got: %+v, want: %+v", got, want)
		}
	})

	t.Run("deletes the payment", func(t *testing.T) {
		err := s.DeletePayment(p2.ID)
		if err != nil {
			t.Errorf("expected no error, got %+v", err)
		}

		newPayments, err := s.GetPayments(store.NewFilter())
		if err != nil {
			t.Fatalf("%+v", err)
		}

		if got, want := newPayments, []*models.Payment{p1, p3, p4}; !reflect.DeepEqual(got, want) {
			t.Errorf("got: %+v, want: %+v", got, want)
		}
	})
}

func TestUpdatePayment(t *testing.T) {
	canonicalProject := paymentFromFixture(t)

	p1 := newPaymentFrom(canonicalProject)
	p1.Attributes.Amount = "0.01"
	p2 := newPaymentFrom(canonicalProject)
	p3 := newPaymentFrom(canonicalProject)
	p4 := newPaymentFrom(canonicalProject)

	payments := []*models.Payment{p1, p2, p3, p4}
	s := store.NewVolatilePaymentsStore(payments)

	t.Run("returns NotFound if not found", func(t *testing.T) {
		p := newPaymentFrom(p1)
		p.ID = strfmt.UUID(uuid.New().String())

		_, err := s.UpdatePayment(p.ID, p)
		if !errors.IsNotFound(err) {
			t.Errorf("expected NotFound, got %+v", err)
		}

		newPayments, err := s.GetPayments(store.NewFilter())
		if err != nil {
			t.Fatalf("%+v", err)
		}

		if got, want := newPayments, payments; !reflect.DeepEqual(got, want) {
			t.Errorf("got: %+v, want: %+v", got, want)
		}
	})

	t.Run("returns BadRequest if request ID does not match embedded ID", func(t *testing.T) {
		p := newPaymentFrom(p1)
		p.ID = strfmt.UUID(uuid.New().String())

		_, err := s.UpdatePayment(strfmt.UUID(uuid.New().String()), p)
		if !errors.IsBadRequest(err) {
			t.Errorf("expected BadRequest, got %+v", err)
		}

		newPayments, err := s.GetPayments(store.NewFilter())
		if err != nil {
			t.Fatalf("%+v", err)
		}

		if got, want := newPayments, payments; !reflect.DeepEqual(got, want) {
			t.Errorf("got: %+v, want: %+v", got, want)
		}
	})

	t.Run("returns BadRequest if validation for Payment fails", func(t *testing.T) {
		p := newPaymentFrom(p1)
		p.ID = strfmt.UUID(uuid.New().String())
		// Type is required to have length > 0
		p.Type = ""

		_, err := s.UpdatePayment(p.ID, p)
		if !errors.IsBadRequest(err) {
			t.Errorf("expected BadRequest, got %+v", err)
		}

		newPayments, err := s.GetPayments(store.NewFilter())
		if err != nil {
			t.Fatalf("%+v", err)
		}

		if got, want := newPayments, payments; !reflect.DeepEqual(got, want) {
			t.Errorf("got: %+v, want: %+v", got, want)
		}
	})

	t.Run("updates a payment", func(t *testing.T) {
		newP3 := clonePayment(p3)
		newP3.Attributes.Amount = "50.00"

		res, err := s.UpdatePayment(newP3.ID, newP3)
		if err != nil {
			t.Errorf("expected no error, got %+v", err)
		}

		newPayments, err := s.GetPayments(store.NewFilter())
		if err != nil {
			t.Fatalf("%+v", err)
		}

		if got, want := res, newP3; !reflect.DeepEqual(got, want) {
			t.Errorf("Return value does not match. got: %+v, want: %+v", got, want)
		}

		if got, want := newPayments, []*models.Payment{p1, p2, newP3, p4}; !reflect.DeepEqual(got, want) {
			t.Errorf("got: %+v, want: %+v", got, want)
		}
	})
}
