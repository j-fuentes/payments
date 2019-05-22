package store

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/j-fuentes/payments/pkg/models"
	"github.com/juju/errors"
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
		if filter.OrganisationID != nil && filter.OrganisationID.String() != p.OrganisationID.String() {
			continue
		}

		amount, err := strconv.ParseFloat(p.Attributes.Amount, 64)
		if err != nil {
			// This should never happen, since we assume data in store in legal
			panic(err)
		}

		if amount > filter.MaxAmount || amount < filter.MinAmount {
			continue
		}

		result = append(result, p)
	}

	return result, nil
}

func (s *VolatilePaymentsStore) GetPayment(id strfmt.UUID) (*models.Payment, error) {
	for _, p := range s.payments {
		if id.String() == p.ID.String() {
			return p, nil
		}
	}

	return nil, errors.NotFoundf("Payment with id %q", id.String())
}
