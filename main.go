package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Course struct {
	CourseNumber  string `json:"course_number"`
	CourseName    string `json:"course_name"`
	Credits       int    `json:"credits"`
	Description   string `json:"description"`
	Prerequisites string `json:"prerequisites"`
	Corequisites  string `json:"corequisites"`
}

type Message struct {
	Message string `json:"message"`
}

// Handler for the "/api/hello" endpoint
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	response := Message{Message: "Hello from the Go backend!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handler for the "/catalog" endpoint to serve course data from the JSON file
func catalogHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Open and read the JSON file
	jsonFile, err := os.Open("data/catalog.json")
	if err != nil {
		http.Error(w, "Unable to open JSON file", http.StatusInternalServerError)
		fmt.Println(err)
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

func main() {
	// Define route handlers
	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/catalog", catalogHandler)

	// Start the server
	log.Println("Backend server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
