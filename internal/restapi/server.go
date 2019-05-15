package restapi

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/j-fuentes/payments/internal/restapi/controllers"
)

// Serve serves the api
func Serve(addr string) error {
	r := mux.NewRouter()

	// Mount routes
	r.HandleFunc("/payments", controllers.GetPayments)

	glog.Infof("Listening on %s", addr)
	return http.ListenAndServe(addr, r)
}
