package restapi

import (
	"net/http"
	"path"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/j-fuentes/payments/internal/restapi/helpers"
	"github.com/j-fuentes/payments/internal/store"
	"github.com/j-fuentes/payments/pkg/models"
	"github.com/juju/errors"
)

func (server *PaymentsServer) GetPayments(w http.ResponseWriter, r *http.Request) {
	filter := store.NewFilter()

	if orgID := r.URL.Query().Get("organisation-id"); orgID != "" {
		aux := strfmt.UUID(orgID)
		filter.OrganisationID = &aux
	}

	if minAmount := r.URL.Query().Get("min-amount"); minAmount != "" {
		min, err := strconv.ParseFloat(minAmount, 64)
		if err != nil {
			helpers.WriteError(w, 422, errors.BadRequestf("invalid format in min-amount parameter"))
		}
		filter.MinAmount = min
	}

	if maxAmount := r.URL.Query().Get("max-amount"); maxAmount != "" {
		max, err := strconv.ParseFloat(maxAmount, 64)
		if err != nil {
			helpers.WriteError(w, 422, errors.BadRequestf("invalid format in max-amount parameter"))
		}
		filter.MaxAmount = max
	}

	payments, err := server.paymentsStore.GetPayments(filter)
	if err != nil {
		// TODO: use a helper for proper error handling
		panic(err)
	}

	result := &models.Payments{
		Data: payments,
		Links: &models.PaymentsLinks{
			Self: helpers.DerefString(path.Join(server.externalURL, r.URL.EscapedPath())),
		},
	}

	helpers.WriteRes(w, result)
}
