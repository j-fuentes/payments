package restapi

import (
	"net/http"

	"github.com/golang/glog"
)

// type loggingResponseWriter struct {
// 	http.ResponseWriter
// 	requestID  string
// 	statusCode int
// }

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		glog.Infof("--> %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func logResponse(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		glog.Infof("<-- %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}
