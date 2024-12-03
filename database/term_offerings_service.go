package database

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

type CourseResult struct {
	Key          string `json:"key"`
	Code         string `json:"code"`
	Title        string `json:"title"`
	Crn          string `json:"crn"`
	No           string `json:"no"`
	Total        string `json:"total"`
	Schd         string `json:"schd"`
	Camp         string `json:"camp"`
	Stat         string `json:"stat"`
	SsrFees      string `json:"ssrFees"`
	IsCancelled  string `json:"isCancelled"`
	Meets        string `json:"meets"`
	Mpkey        string `json:"mpkey"`
	MeetingTimes string `json:"meetingTimes"`
	Instr        string `json:"instr"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Srcdb        string `json:"srcdb"`
	Term         string `json:"term"`
}

type CourseOffering struct {
	Count   int            `json:"count"`
	Results []CourseResult `json:"results"`
	Term    string         `json:"term"`
	Code    string         `json:"code"`
	Title   string         `json:"title"`
}

type TermOffering struct {
	Department string           `json:"department"`
	Term       string           `json:"term"`
	Courses    []CourseOffering `json:"courses"`
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
	w.WriteHeader(http.StatusOK)
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
		filter["courses.code"] = bson.M{"$regex": courseQuery, "$options": "i"}
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
