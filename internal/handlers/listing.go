package handlers

import (
	"net/http"
)

func List(w http.ResponseWriter,r *http.Request){
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(`{"message": "List of products"}`)) 

}