package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/AdviseU-Project/AdviseU-backend/database"
	"github.com/lmittmann/tint"
)

func main() {
	// Set up logging
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, nil)))

	// Initialize MongoDB
	if err := database.InitMongoDB(); err != nil {
		log.Fatalf("Error initializing MongoDB: %v", err)
	}
	defer func() {
		if err := database.MongoClient.Disconnect(database.Ctx); err != nil {
			log.Fatalf("Error disconnecting MongoDB client: %v", err)
		}
	}()

	// Define route handlers
	http.HandleFunc("/catalog", database.CatalogHandler)
	http.HandleFunc("/catalogs", database.CatalogsHandler)
	http.HandleFunc("/term_offerings", database.TermOfferingsHandler)

	// Start the server
	log.Println("backend server is starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", requestLoggingMiddleware(corsMiddleware(http.DefaultServeMux))))
}
