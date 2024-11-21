package main

import (
	"encoding/json"
	"io"
	"io/fs"
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

var availableCatalogs []string
var catalogMap map[string][]Course

func loadCatalogData() error {
	catalogDir := "data/catalogs/"
	availableCatalogs = make([]string, 0)
	catalogMap = make(map[string][]Course)

	err := filepath.WalkDir(catalogDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Check if the file follows the "{department}_catalog.json" naming convention
		if strings.HasSuffix(d.Name(), "_catalog.json") {
			// Extract the department code from the filename
			department := strings.TrimSuffix(d.Name(), "_catalog.json")
			availableCatalogs = append(availableCatalogs, department)

			jsonFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer jsonFile.Close()

			byteValue, err := io.ReadAll(jsonFile)
			if err != nil {
				return err
			}

			catalog := make([]Course, 0)
			err = json.Unmarshal(byteValue, &catalog)
			if err != nil {
				return err
			}

			catalogMap[department] = catalog
		}

		return nil
	})

	return err
}

func getListOfCatalogs() ([]string, error) {
	return availableCatalogs, nil
}

func queryCoursesFromCatalog(catalogId string, courseQuery string) ([]Course, error) {
	result := make([]Course, 0)

	for catalog, courses := range catalogMap {
		if catalogId != "" && catalog != catalogId {
			continue
		}

		for _, course := range courses {
			if strings.Contains(course.CourseNumber, courseQuery) || strings.Contains(course.CourseName, courseQuery) {
				result = append(result, course)
			}
		}
	}

	return result, nil
}
