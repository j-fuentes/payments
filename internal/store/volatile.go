package store

import (
	"fmt"
	"strconv"

	"github.com/j-fuentes/payments/pkg/models"
)

// VolatilePaymentsStore is an in-memory store for payments. It implements the PaymentsStore interface.
type VolatilePaymentsStore struct {
	payments []*models.Payment
}

// NewVolatilePaymentsStore creates a new in-memory store and initializes it.
func NewVolatilePaymentsStore(payments []*models.Payment) *VolatilePaymentsStore {
	return &VolatilePaymentsStore{
		payments: payments,
	}
}

// GetPayments returns the list of payments that match a given filter.
func (s *VolatilePaymentsStore) GetPayments(filter Filter) ([]*models.Payment, error) {
	result := []*models.Payment{}

	// check if payment satisfies the filter
	for _, p := range s.payments {
		fmt.Println("skipping?")
		if filter.OrganisationID != nil && filter.OrganisationID.String() != p.OrganisationID.String() {
			fmt.Println("skipped")
			continue
		}

		amount, err := strconv.ParseFloat(p.Attributes.Amount, 64)
		if err != nil {
			// This should never happen, since we assume data in store in legal
			panic(err)
		}

		if amount > filter.MaxAmount || amount < filter.MinAmount {
			fmt.Printf("skipped_amount: %f, %f\n", amount, filter.MaxAmount)
			continue
		}

		fmt.Println("madeit")
		result = append(result, p)
	}

	return result, nil
}
