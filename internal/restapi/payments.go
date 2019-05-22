package restapi

import (
	"net/http"
	"path"
	"strconv"
	"io/ioutil"

	"github.com/gorilla/mux"
	"github.com/go-openapi/strfmt"
	"github.com/j-fuentes/payments/internal/restapi/helpers"
	"github.com/j-fuentes/payments/internal/store"
	"github.com/j-fuentes/payments/pkg/models"
	"github.com/juju/errors"
	"github.com/golang/glog"
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
			helpers.WriteError(w, http.StatusBadRequest, errors.BadRequestf("invalid format in min-amount parameter"))
		}
		filter.MinAmount = min
	}

	if maxAmount := r.URL.Query().Get("max-amount"); maxAmount != "" {
		max, err := strconv.ParseFloat(maxAmount, 64)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, errors.BadRequestf("invalid format in max-amount parameter"))
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

func (server *PaymentsServer) GetPayment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := strfmt.UUID(params["id"])

	p, err := server.paymentsStore.GetPayment(id)
	if err != nil {
		if errors.IsNotFound(err) {
			helpers.WriteError(w, http.StatusNotFound, errors.NotFoundf("Cannot find payment with id %q", id))
		} else {
			glog.Errorf("%+v", err)
			helpers.WriteError(w, http.StatusInternalServerError, errors.Errorf("Internal error"))
		}
		return
	}

	helpers.WriteRes(w, p)
}

func (server *PaymentsServer) DeletePayment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := strfmt.UUID(params["id"])

	err := server.paymentsStore.DeletePayment(id)
	if err != nil {
		if errors.IsNotFound(err) {
			helpers.WriteError(w, http.StatusNotFound, errors.NotFoundf("Cannot find payment with id %q", id))
		} else {
			glog.Errorf("%+v", err)
			helpers.WriteError(w, http.StatusInternalServerError, errors.Errorf("Internal error"))
		}
		return
	}

	helpers.WriteRes(w, &models.Empty{})
}

func (server *PaymentsServer) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := strfmt.UUID(params["id"])

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var newPayment models.Payment
	err = newPayment.UnmarshalBinary(b)
	if err != nil {
		glog.Errorf("%+v", err)
		helpers.WriteError(w, http.StatusBadRequest, errors.BadRequestf("Cannot unmarshal payload"))
		return
	}
	if err = newPayment.Validate(nil); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.BadRequestf("Cannot validate payload: %v", err))
		return
	}


	p, err := server.paymentsStore.UpdatePayment(id, &newPayment)
	if err != nil {
		if errors.IsNotFound(err) {
			helpers.WriteError(w, http.StatusNotFound, errors.NotFoundf("Cannot find payment with id %q", id))
		} else if errors.IsBadRequest(err) {
			glog.Errorf("%+v", err)
			helpers.WriteError(w, http.StatusBadRequest, err)
		} else {
			glog.Errorf("%+v", err)
			helpers.WriteError(w, http.StatusInternalServerError, errors.Errorf("Internal error"))
		}
		return
	}

	helpers.WriteRes(w, p)
}
