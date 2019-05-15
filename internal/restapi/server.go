package restapi

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/j-fuentes/payments/internal/restapi/controllers"
)

// middleware is a wrapper funcion that performs an arbitrary operation with the incoming request and then calls another htt.HandlerFunc. Ideally middlewares can be chained.
type middleware func(next http.HandlerFunc) http.HandlerFunc

// chainMiddleware returna a handler as a result of chaining the ones received as parameters.
func chainMiddleware(mw ...middleware) middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last(w, r)
		}
	}
}

// Serve serves the api
func Serve(addr string) error {
	handle := chainMiddleware(logRequest, logResponse)

	r := mux.NewRouter()

	// Mount routes
	r.HandleFunc("/payments", handle(controllers.GetPayments))

	glog.Infof("Listening on %s", addr)
	return http.ListenAndServe(addr, r)
}
