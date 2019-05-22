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

func (s *VolatilePaymentsStore) CreatePayment(newPayment *models.Payment) (*models.Payment, error) {
	if err := newPayment.Validate(nil); err != nil {
		return nil, errors.BadRequestf("Payment not valid", err)
	}

	s.payments = append(s.payments, newPayment)

	return newPayment, nil
}

func (s *VolatilePaymentsStore) GetPayment(id strfmt.UUID) (*models.Payment, error) {
	for _, p := range s.payments {
		if id.String() == p.ID.String() {
			return p, nil
		}
	}

	return nil, errors.NotFoundf("Payment with id %q", id.String())
}

func (s *VolatilePaymentsStore) DeletePayment(id strfmt.UUID) error {
	for i, p := range s.payments {
		if id.String() == p.ID.String() {
			s.payments = append(s.payments[:i], s.payments[i+1:]...)
			return nil
		}
	}

	return errors.NotFoundf("Payment with id %q", id.String())
}

func (s *VolatilePaymentsStore) UpdatePayment(id strfmt.UUID, newPayment *models.Payment) (*models.Payment, error) {
	if k, e := id.String(), newPayment.ID.String(); k != e {
		return nil, errors.BadRequestf("Provided ID (%q) does not match embedded ID in new Payment (%q)", k, e)
	}

	if err := newPayment.Validate(nil); err != nil {
		return nil, errors.BadRequestf("Payment not valid", err)
	}

	for i, p := range s.payments {
		if id.String() == p.ID.String() {
			s.payments[i] = newPayment
			return newPayment, nil
		}
	}

	return nil, errors.NotFoundf("Payment with id %q", id.String())
}
