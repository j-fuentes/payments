package store

import (
	"reflect"
	"testing"

	"github.com/j-fuentes/payments/pkg/models"
)

func TestGetPayments(t *testing.T) {
	genDesc := "hello world"
	payments := []*models.Payment{
		&models.Payment{ID: 1, Description: &genDesc},
		&models.Payment{ID: 2, Description: &genDesc},
	}

	testCases := []struct {
		name         string
		filter       Filter
		wantPayments []*models.Payment
	}{
		{
			"returns all the payments if empty filter",
			Filter{},
			payments,
		},
	}

	s := VolatilePaymentsStore{
		payments: payments,
	}

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
