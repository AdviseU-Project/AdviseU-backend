package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/lmittmann/tint"
)

func main() {
	// Set up logging
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, nil)))

	// Load backing course data
	err := loadCatalogData()
	if err != nil {
		slog.With("err", err).Error("Could not load catalog data from disk")
		return
	}

	// Define route handlers
	http.HandleFunc("/catalog", catalogHandler)
	http.HandleFunc("/catalogs", catalogsHandler)

	// Start the server
	log.Println("backend server is starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", requestLoggingMiddleware(corsMiddleware(http.DefaultServeMux))))
}
