package httpx

import (
	"encoding/json"
	"net/http"
)

func Error(w http.ResponseWriter, msg string, code string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": msg,
		"code":  code,
	})
}