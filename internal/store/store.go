package store

import (
	"math"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/j-fuentes/payments/pkg/models"
)

// Filter contains attributes payments can be filtered by.
type Filter struct {
	OrganisationID *strfmt.UUID
	MaxAmount      float64
	MinAmount      float64
	// this can be expanded to filter by any desired criteria, e.g. currency, sender, debtor, etc.
}

// NewFilter returns an all-pass filter that can be edited to add restrictions.
func NewFilter() Filter {
	return Filter{
		OrganisationID: nil,
		MaxAmount:      math.Inf(1),
		MinAmount:      0,
	}
}

// PaymentsStore defines the interface of a store of payments.
type PaymentsStore interface {
	// GetPayments returns the list of payments that match a given filter.
	GetPayments(filter Filter) ([]*models.Payment, error)
	// CreatePayment creates a new Payment.
	CreatePayment(newPayment *models.Payment) (*models.Payment, error)
	// GetPayment returns a Payment by its ID.
	GetPayment(id strfmt.UUID) (*models.Payment, error)
	// DeletePayment deletes a Payment by its ID.
	DeletePayment(id strfmt.UUID) error
	// UpdatePayment updates a Payment by its ID.
	UpdatePayment(id strfmt.UUID, newPayment *models.Payment) (*models.Payment, error)
}
