package restapi

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/j-fuentes/payments/internal/store"
	"github.com/j-fuentes/payments/internal/restapi/helpers"
	"github.com/juju/errors"
)

// PaymentsServer serves a REST API for payments.
type PaymentsServer struct {
	paymentsStore store.PaymentsStore
	externalURL   string
}

// NewPaymentsServer creates a new PaymentsServer that consumes data from a PaymentsStore.
func NewPaymentsServer(s store.PaymentsStore, externalURL string) *PaymentsServer {
	return &PaymentsServer{
		paymentsStore: s,
		externalURL:   externalURL,
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	helpers.WriteError(w, 404, errors.NotFoundf("The requested resource does not exist"))
}

func withJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}
}

// Serve serves the REST API.
func (server *PaymentsServer) Serve(addr string) error {
	handle := chainMiddleware(withJSON, withRequestID, withLogging)

	r := mux.NewRouter()

	// Mount routes
	r.HandleFunc("/payments", handle(server.GetPayments))
	r.NotFoundHandler = handle(http.HandlerFunc(notFoundHandler))

	glog.Infof("Listening on %s", addr)
	return http.ListenAndServe(addr, r)
}
