package database

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

type Course struct {
	Department    string     `json:"department" bson:"department"`
	CourseNumber  string     `json:"course_number" bson:"course_number"`
	CourseName    string     `json:"course_name" bson:"course_name"`
	Credits       string     `json:"credits" bson:"credits"`
	Description   string     `json:"description" bson:"description"`
	Prerequisites [][]string `json:"prerequisites" bson:"prerequisites"`
	Corequisites  [][]string `json:"corequisites" bson:"corequisites"`
}

// Handler for the "/catalog" endpoint to serve course data from the database
func CatalogHandler(w http.ResponseWriter, r *http.Request) {
	// Get the request from query parameters
	catalogId := r.URL.Query().Get("catalog_id")
	courseQuery := r.URL.Query().Get("course_query")

	slog.With("catalog_id", catalogId, "course_query", courseQuery).Info("parameters")

	// Query MongoDB for courses
	courses, err := QueryCoursesFromCatalog(catalogId, courseQuery)
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
func CatalogsHandler(w http.ResponseWriter, r *http.Request) {
	// Define the directory where catalog files are stored
	availableCatalogs, err := GetListOfCatalogs()
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

// QueryCoursesFromCatalog fetches courses from the MongoDB `catalogs` collection
func QueryCoursesFromCatalog(catalogId string, courseQuery string) ([]Course, error) {
	collection := MongoClient.Database("adviseu_db").Collection("catalogs")

	// Build query filter
	filter := bson.M{}
	if catalogId != "" {
		filter["department"] = catalogId
	}
	if courseQuery != "" {
		filter["$or"] = []bson.M{
			{"course_number": bson.M{"$regex": courseQuery, "$options": "i"}},
			{"course_name": bson.M{"$regex": courseQuery, "$options": "i"}},
		}
	}

	// Query MongoDB
	cursor, err := collection.Find(Ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(Ctx)

	// Decode results
	var results []Course
	if err = cursor.All(Ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// GetListOfCatalogs retrieves distinct department codes from MongoDB
func GetListOfCatalogs() ([]string, error) {
	collection := MongoClient.Database("adviseu_db").Collection("catalogs")

	// Use MongoDB distinct query to get all departments
	departments, err := collection.Distinct(Ctx, "department", bson.M{})
	if err != nil {
		return nil, err
	}

	// Convert results to a list of strings
	var result []string
	for _, dept := range departments {
		result = append(result, dept.(string))
	}

	return result, nil
}
