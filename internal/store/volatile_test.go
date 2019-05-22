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

	orgID := strfmt.UUID(uuid.New().String())
	minAmount := 100.2
	maxAmount := 300.0

	canonicalProject := fixture.Data[0]

	p1 := copyPayment(canonicalProject)

	p2 := copyPayment(canonicalProject)
	p2.OrganisationID = orgID

	p3 := copyPayment(canonicalProject)
	p3.Attributes.Amount = fmt.Sprintf("%.2f", minAmount+1.0)

	p4 := copyPayment(canonicalProject)
	p4.OrganisationID = orgID
	p4.Attributes.Amount = fmt.Sprintf("%.2f", minAmount+1.0)

	p5 := copyPayment(canonicalProject)
	p5.OrganisationID = orgID
	p5.Attributes.Amount = fmt.Sprintf("%.2f", minAmount-1.0)

	p6 := copyPayment(canonicalProject)
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

func TestGetPayment(t *testing.T) {
	file := "single.json"
	fixture, err := fixtures.LoadPayments(file)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if l := len(fixture.Data); l != 1 {
		t.Fatalf("expected to load just one payment from %s, but found %d", file, l)
	}

	canonicalProject := fixture.Data[0]

	p1 := copyPayment(canonicalProject)
	p2 := copyPayment(canonicalProject)

	s := store.NewVolatilePaymentsStore(
		[]*models.Payment{p1, p2},
	)

	t.Run("returns 404 if not found", func(t *testing.T) {
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
