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
	"github.com/google/uuid"
)

// GetPayments writes a list of Payments filtered by 'organisation-id', 'min-amount' and 'max-amount'.
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
			return
		}
		filter.MinAmount = min
	}

	if maxAmount := r.URL.Query().Get("max-amount"); maxAmount != "" {
		max, err := strconv.ParseFloat(maxAmount, 64)
		if err != nil {
			helpers.WriteError(w, http.StatusBadRequest, errors.BadRequestf("invalid format in max-amount parameter"))
			return
		}
		filter.MaxAmount = max
	}

	payments, err := server.paymentsStore.GetPayments(filter)
	if err != nil {
		glog.Errorf("Cannot GetPayments: %+v", err)
		helpers.WriteError(w, http.StatusInternalServerError, errors.Errorf("Cannot retrieve list of payments."))
		return
	}

	result := &models.Payments{
		Data: payments,
		Links: &models.PaymentsLinks{
			Self: helpers.DerefString(path.Join(server.externalURL, r.URL.EscapedPath())),
		},
	}

	helpers.WriteRes(w, result)
}

// CreatePayment creates a new Payment with a new ID.
func (server *PaymentsServer) CreatePayment(w http.ResponseWriter, r *http.Request) {
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

	newPayment.ID = strfmt.UUID(uuid.New().String())

	if err = newPayment.Validate(nil); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.BadRequestf("Cannot validate payload: %v", err))
		return
	}


	p, err := server.paymentsStore.CreatePayment(&newPayment)
	if err != nil {
		if errors.IsBadRequest(err) {
			glog.Errorf("%+v", err)
			helpers.WriteError(w, http.StatusBadRequest, err)
		} else {
			glog.Errorf("%+v", err)
			helpers.WriteError(w, http.StatusInternalServerError, errors.Errorf("Internal error"))
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	helpers.WriteRes(w, p)
}

// GetPayment writes a Payment with a given 'id'.
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

// DeletePayment deletes a payment with a given 'id'.
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

// UpdatePayment updates a payment with a given 'id'.
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
