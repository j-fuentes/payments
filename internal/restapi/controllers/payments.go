package controllers

import (
	"net/http"

	"github.com/j-fuentes/payments/internal/restapi/helpers"
	"github.com/j-fuentes/payments/pkg/models"
)

func GetPayments(w http.ResponseWriter, r *http.Request) {
	genDesc := "hello world"
	payments := &models.Payments{
		Data: []*models.Payment{
			&models.Payment{ID: 1, Description: &genDesc},
			&models.Payment{ID: 2, Description: &genDesc},
		},
		Links: &models.PaymentsLinks{
			Self: helpers.DerefString(helpers.GenerateLink(r.URL.EscapedPath())),
		},
	}

	helpers.WriteRes(w, payments)
}
