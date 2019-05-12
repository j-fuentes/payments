package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/j-fuentes/payments/pkg/models"
)

func GetPayments(w http.ResponseWriter, r *http.Request) {
	payments := []models.Payment{
		models.Payment{ID: 1},
		models.Payment{ID: 2},
	}

	// validate against spec
	json.NewEncoder(w).Encode(payments)
}
