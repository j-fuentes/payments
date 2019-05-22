package store_test

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/getlantern/deepcopy"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/j-fuentes/payments/internal/fixtures"
	"github.com/j-fuentes/payments/internal/store"
	"github.com/j-fuentes/payments/pkg/models"
)

func TestGetPayments(t *testing.T) {
	file := "single.json"
	fixture, err := fixtures.LoadPayments(file)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if l := len(fixture.Data); l != 1 {
		t.Fatalf("expected to load just one payment from %s, but found %d", file, l)
	}

	orgID := strfmt.UUID("aaaaaaaa-bbbb-cccc-dddd-ddeadbeefc43")
	minAmount := 100.2
	maxAmount := 300.0

	canonicalProject := fixture.Data[0]

	var p1 models.Payment
	deepcopy.Copy(&p1, canonicalProject)

	var p2 models.Payment
	deepcopy.Copy(&p2, canonicalProject)
	p2.OrganisationID = orgID

	var p3 models.Payment
	deepcopy.Copy(&p3, canonicalProject)
	p3.Attributes.Amount = fmt.Sprintf("%.2f", minAmount+1.0)

	var p4 models.Payment
	deepcopy.Copy(&p4, canonicalProject)
	p4.OrganisationID = orgID
	p4.Attributes.Amount = fmt.Sprintf("%.2f", minAmount+1.0)

	var p5 models.Payment
	deepcopy.Copy(&p5, canonicalProject)
	p5.OrganisationID = orgID
	p5.Attributes.Amount = fmt.Sprintf("%.2f", minAmount-1.0)

	var p6 models.Payment
	deepcopy.Copy(&p6, canonicalProject)
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
			[]*models.Payment{&p1, &p2, &p3, &p4, &p5, &p6},
		},
		{
			"returns all the payments of one organisation",
			store.Filter{
				OrganisationID: &orgID,
				MaxAmount:      math.Inf(1),
				MinAmount:      0,
			},
			[]*models.Payment{&p2, &p4, &p5, &p6},
		},
		{
			"returns all the payments between MinAmount and MaxAmount",
			store.Filter{
				OrganisationID: nil,
				MinAmount:      minAmount + 2.0,
				MaxAmount:      maxAmount + 2.0,
			},
			[]*models.Payment{&p6},
		},
	}

	s := store.NewVolatilePaymentsStore(
		[]*models.Payment{&p1, &p2, &p3, &p4, &p5, &p6},
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
