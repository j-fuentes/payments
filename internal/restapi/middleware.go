package restapi

import "net/http"

// middleware is a wrapper funcion that performs an arbitrary operation with the incoming request and then calls another htt.HandlerFunc. Ideally middlewares can be chained.
type middleware func(next http.HandlerFunc) http.HandlerFunc

// chainMiddleware returns a handler as a result of chaining the ones received as parameters.
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
