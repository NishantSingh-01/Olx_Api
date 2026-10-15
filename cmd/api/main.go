package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"olx_api/internal/config"
	"olx_api/internal/database"
	"olx_api/internal/handlers"
)

func main() {
	cfg := config.MustLoad()
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("main: database.Connect: Database Connection failed : %v", err)
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	fmt.Println("Server is starting on port", cfg.Port)

	listingHandler := handlers.NewListing_Handler(db,logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.HealthzHandler)
	//listing handler
	mux.HandleFunc("GET /listing", listingHandler.ListingHandler)
	mux.HandleFunc("DELETE /listing/{id}", listingHandler.DeleteListing)
	mux.HandleFunc("POST /listing", listingHandler.CreateListing)
	mux.HandleFunc("GET /listing/{id}", listingHandler.GetListingByID)
	mux.HandleFunc("PUT /listing/{id}", handlers.UpdateListing(db))

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
