package store

import (
	"github.com/j-fuentes/payments/pkg/models"
)

// Filter contains attributes payments can be filtered by.
type Filter struct {
	// TODO: add actual fields to filter by
}

// PaymentsStore defines the interface of a store of payments.
type PaymentsStore interface {
	GetPayments(filter Filter) ([]*models.Payment, error)
}
