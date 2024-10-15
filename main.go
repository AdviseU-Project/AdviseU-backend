package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Course struct {
	CourseNumber  string     `json:"course_number"`
	CourseName    string     `json:"course_name"`
	Credits       string     `json:"credits"`
	Description   string     `json:"description"`
	Prerequisites [][]string `json:"prerequisites"`
	Corequisites  [][]string `json:"corequisites"`
}

// Handler for the "/catalog" endpoint to serve course data from the JSON file
func catalogHandler(w http.ResponseWriter, r *http.Request) {
	// Get the department from query parameters, default to "CS" if not provided
	department := r.URL.Query().Get("department")

	// Construct the file path based on department (e.g., "data/CS_catalog.json")
	filePath := fmt.Sprintf("data/catalogs/%s_catalog.json", department)

	// Try to open the corresponding JSON file
	jsonFile, err := os.Open(filePath)
	if err != nil || department == "" {
		// If file doesn't exist or dept is empty, return an empty course list
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Course{}) // Return an empty slice of Course
		return
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var catalog []Course
	json.Unmarshal(byteValue, &catalog)

	// Serve JSON as HTTP response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(catalog)
}

// Handler for the "/catalogs" endpoint to return the list of available catalogs
func catalogsHandler(w http.ResponseWriter, r *http.Request) {
	// Define the directory where catalog files are stored
	catalogDir := "data/catalogs/"

	// Find all files that end with "_catalog.json"
	var availableCatalogs []string
	err := filepath.WalkDir(catalogDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Check if the file follows the "{department}_catalog.json" naming convention
		if strings.HasSuffix(d.Name(), "_catalog.json") {
			// Extract the department code from the filename
			department := strings.TrimSuffix(d.Name(), "_catalog.json")
			availableCatalogs = append(availableCatalogs, department)
		}
		return nil
	})

	if err != nil {
		http.Error(w, "Error reading catalog directory", http.StatusInternalServerError)
		return
	}

	// Serve the list of available catalogs as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(availableCatalogs)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Define route handlers
	http.HandleFunc("/catalog", catalogHandler)
	http.HandleFunc("/catalogs", catalogsHandler)

	// Start the server
	log.Println("Backend server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(http.DefaultServeMux)))
}
