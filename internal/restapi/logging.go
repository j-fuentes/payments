package restapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/google/uuid"
)

type contextKey string

const requestIDKey = contextKey("requestID")

func getRequestIDFromContext(ctx context.Context) uuid.UUID {
	reqID, ok := ctx.Value(requestIDKey).(uuid.UUID)
	if !ok {
		panic("UUID not found in request context")
	}
	return reqID
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

const statusKey = contextKey("statusCode")

func getStatusFromContext(ctx context.Context) int {
	status, ok := ctx.Value(statusKey).(int)
	if !ok {
		panic("Status code not found in request context")
	}
	return status
}

// Middlewares

func withRequestID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("withRequestID")
		ctx := r.Context()
		ctx = context.WithValue(ctx, requestIDKey, uuid.New())
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Initialize the status to 200 in case WriteHeader is not called
		rec := statusRecorder{w, 200}

		next.ServeHTTP(&rec, r)
		glog.Infof("- %q %s %s -> %d", getRequestIDFromContext(r.Context()), r.Method, r.URL.Path, rec.status)
	}
}
