package database

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

type TermOffering struct {
}

// TermOfferingsHandler handles requests for term offerings
func TermOfferingsHandler(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	courseQuery := r.URL.Query().Get("course_query")

	offerings, err := QueryTermOfferings(term, courseQuery)
	if err != nil {
		http.Error(w, "Error querying term offerings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(offerings)
}

// QueryTermOfferings fetches offerings from the MongoDB `term_offerings` collection
func QueryTermOfferings(term, courseQuery string) ([]TermOffering, error) {
	collection := MongoClient.Database("adviseu_db").Collection("term_offerings")
	filter := bson.M{}

	if term != "" {
		filter["term"] = term
	}
	if courseQuery != "" {
		filter["course_number"] = bson.M{"$regex": courseQuery, "$options": "i"}
	}

	cursor, err := collection.Find(Ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(Ctx)

	var results []TermOffering
	if err := cursor.All(Ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
