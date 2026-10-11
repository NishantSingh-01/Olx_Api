package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"olx_api/internal/config"
	"olx_api/internal/database"
	"olx_api/internal/handlers"
)

func main() {
	cfg := config.MustLoad()
	_, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("main: database.Connect: Database Connection failed : %v", err)
	}
	

	fmt.Println("Server is starting on port", cfg.Port)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health",handlers.HealthzHandler)

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
