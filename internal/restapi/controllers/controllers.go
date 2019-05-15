package controllers

import (
	"net/http"

	"github.com/j-fuentes/payments/internal/restapi/helpers"
	"github.com/j-fuentes/payments/pkg/models"
)

func GetPayments(w http.ResponseWriter, r *http.Request) {
	payments := &models.Payments{
		Data: []*models.Payment{
			&models.Payment{ID: 1},
			&models.Payment{ID: 2},
		},
		Links: &models.PaymentsLinks{},
	}

	helpers.WriteRes(w, payments)
}
