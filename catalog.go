package main

// Handler for the "/catalog" endpoint to serve course data from the JSON file
func catalogHandler(w http.ResponseWriter, r *http.Request) {
	// Get the department from query parameters, default to "CS" if not provided
	catalogId := r.URL.Query().Get("catalog_id")
	courseQuery := r.URL.Query().Get("course_query")

	courses, err := queryCoursesFromCatalog(catalogId, courseQuery)
	if err != nil {
		slog.With("err", err,
					"method", r.Method,
					"query", r.URL.RawQuery,
					"client", r.RemoteAddr,
					"url", r.URL.Path).Error("error querying course catalog")
		http.Error(w, "error querying course catalog", http.StatusInternalServerError)
		return
	}

	// Serve JSON as HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(courses)
}

// Handler for the "/catalogs" endpoint to return the list of available catalogs
func catalogsHandler(w http.ResponseWriter, r *http.Request) {
	// Define the directory where catalog files are stored
	availableCatalogs, err := getListOfCatalogs()
	if err != nil {
		slog.With("err", err,
					"method", r.Method,
					"client", r.RemoteAddr,
					"url", r.URL.Path).Error("error reading catalog directory")
		http.Error(w, "error reading catalog directory", http.StatusInternalServerError)
		return
	}

	// Serve the list of available catalogs as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(availableCatalogs)
}