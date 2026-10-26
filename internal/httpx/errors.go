package httpx

import (
	"encoding/json"
	"net/http"
)

type Code string

const (
	Errinvalid_id        Code = "invalid_id"
	Errnot_found         Code = "not_found"
	Errinternal_error    Code = "internal_error"
	Errmalformed_json    Code = "malformed_json"
	Errvalidation_failed Code = "validation_failed"
	Errunauthenticated   Code = "unauthenticated"
	Errforbidden         Code = "forbidden"
	Errconflict          Code = "conflict"
	Errrate_limited      Code = "rate_limited"
	ErrUnauthorized      Code = "unauthorized"
)

type ErrorEnvelope struct {
	Error ErrorPayLoad `json:"error"`
}
type ErrorPayLoad struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

func Error(w http.ResponseWriter, msg string, code Code, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(
		ErrorEnvelope{
			Error: ErrorPayLoad{
				Message: msg,
				Code:    code,
			},
		},
	)
}
func ValidationError(w http.ResponseWriter, msg string, code Code, field string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(
		ErrorEnvelope{
			Error: ErrorPayLoad{
				Message: msg,
				Code:    code,
				Field:   field,
			},
		},
	)
}

// HTTP Status Code    Custom Code
// 400                 invalid_id
// 404                 not_found
// 500                 internal_error
// 400                 malformed_json
// 422                 validation_failed
// 401                 unauthenticated
// 403                 forbidden
// 409                 conflict
// 429                 rate_limited
