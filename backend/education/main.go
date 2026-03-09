package main

import (
	"log"
	"net/http"

	educationv1 "github.com/daniellawrence/cv/gen/go/education/v1"

	"google.golang.org/protobuf/encoding/protojson"
)

var educations = []*educationv1.Education{
	{
		Id:          1,
		Institution: "Griffith University",
		Degree:      "Bachelor, Information Technology",
		StartDate:   "2005-12",
		EndDate:     "2008-06",
	},
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func listEducation(w http.ResponseWriter, r *http.Request) {
	m := protojson.MarshalOptions{
		UseProtoNames:   false,
		EmitUnpopulated: true,
	}

	w.Header().Set("Content-Type", "application/json")

	data, err := m.Marshal(&educationv1.ListEducationResponse{
		Education: educations,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/education", listEducation)

	log.Println("education service listening on :8080")

	err := http.ListenAndServe(":8080", cors(mux))
	if err != nil {
		log.Fatal(err)
	}
}
