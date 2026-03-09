package main

import (
	"encoding/json"
	"log"
	"net/http"

	educationv1 "github.com/daniellawrence/cv/gen/go/education/v1"
)

var educations = []*educationv1.Education{
	{
		Id:          "1",
		Institution: "Griffith University",
		Degree:      "Bachelor, Information Technology",
		StartDate:   "2005-12",
		EndDate:     "2008-06",
	},
}

func listEducation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(educations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/education", listEducation)

	log.Println("education service listening on :8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
