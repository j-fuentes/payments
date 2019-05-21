package restapi

import (
	"context"
	"net/http"

	"github.com/golang/glog"
	"github.com/google/uuid"
)

type contextKey string

const requestIDKey = contextKey("requestID")

func getRequestIDFromContext(ctx context.Context) uuid.UUID {
	reqID, ok := ctx.Value(requestIDKey).(uuid.UUID)
	if !ok {
		panic("UUID not found or not correct in request context")
	}
	return reqID
}

func withRequestID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, requestIDKey, uuid.New())
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func withRequestLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		glog.Infof("--> ReqID:%q %s %s", getRequestIDFromContext(r.Context()), r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func withResponseLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		glog.Infof("<-- ReqID:%q", getRequestIDFromContext(r.Context()))
		next.ServeHTTP(w, r)
	}
}
