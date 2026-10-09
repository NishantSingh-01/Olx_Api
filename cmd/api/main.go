package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"olx_api/internal/config"
	"time"
)

func main() {
	cfg:= config.MustLoad()

	fmt.Println("Server is starting on port", cfg.Port)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json") //for setting up response header
		w.WriteHeader(http.StatusOK)                       // to set status code 200

		w.Write([]byte(`{
			"status": "ok",
			"message": "server is Healthy"
		}`))
	})

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	// err := http.ListenAndServe(":8090", mux)  // with default time of 0
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server failed to start on port %s %v", cfg.Port, err)
	}
}
