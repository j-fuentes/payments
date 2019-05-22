package restapi

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/j-fuentes/payments/internal/store"
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

// Serve serves the REST API.
func (server *PaymentsServer) Serve(addr string) error {
	handle := chainMiddleware(withRequestID, withLogging)

	r := mux.NewRouter()

	// Mount routes
	r.HandleFunc("/payments", handle(server.GetPayments))

	glog.Infof("Listening on %s", addr)
	return http.ListenAndServe(addr, r)
}
