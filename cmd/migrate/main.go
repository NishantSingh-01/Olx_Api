package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const (
	requestIDKey = "X-Request-ID"
)

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(requestIDKey)
		if id == "" {
			id = uuid.New().String()
		}
		w.Header().Add(requestIDKey, id)
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
