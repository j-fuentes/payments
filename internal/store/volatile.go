package store

import (
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
	// TODO: filter by 'filter' when it is defined
	return s.payments, nil
}
