package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey int 
const (
	requestIdKey ctxKey = iota
)

const (
	requestIDHeader = "X-Request-ID"
)

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestIDHeader)
		if id == "" {
			id = uuid.NewString()
		}
		w.Header().Add(requestIDHeader, id)
		ctx := context.WithValue(r.Context(), requestIdKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(ctx context.Context) string {
	id, ok := ctx.Value(requestIdKey).(string)
	if !ok {
		return ""
	}
	return id
}
