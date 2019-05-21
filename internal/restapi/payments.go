package restapi

import (
	"net/http"

	"github.com/j-fuentes/payments/internal/restapi/helpers"
	"github.com/j-fuentes/payments/internal/store"
	"github.com/j-fuentes/payments/pkg/models"
)

func (server *PaymentsServer) GetPayments(w http.ResponseWriter, r *http.Request) {
	payments, err := server.paymentsStore.GetPayments(store.Filter{})
	if err != nil {
		// TODO: use a helper for proper error handling
		panic(err)
	}

	result := &models.Payments{
		Data: payments,
		Links: &models.PaymentsLinks{
			Self: helpers.DerefString(helpers.GenerateLink(r.URL.EscapedPath())),
		},
	}

	helpers.WriteRes(w, result)
}
