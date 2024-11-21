package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	log.Println("Backend server is starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(http.DefaultServeMux)))
}
