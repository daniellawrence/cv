package main

import (
	"log"
	"net/http"

	"github.com/daniellawrence/cv/backend/common"
	interestv1 "github.com/daniellawrence/cv/gen/go/interest/v1"

	"google.golang.org/protobuf/encoding/protojson"
)

var interests = []*interestv1.Interest{
	{
		Id:   "tech",
		Type: "Technical Interests",
		Names: []string{
			"Site Reliability Engineering",
			"Observability Architecture",
			"Continuous Deployment",
			"Developer Tooling",
		},
	},
	{
		Id:   "skills",
		Type: "Languages/Skills",
		Names: []string{
			"Site Reliability Engineering",
			"Observability Architecture",
			"Continuous Deployment",
			"Developer Tooling",
		},
	},
	{
		Id:   "hobbies",
		Type: "hobbies",
		Names: []string{
			"3D printing / All things maker",
			"Home Automation",
			"Wildlife Rescue",
		},
	},
}

func listInterest(w http.ResponseWriter, r *http.Request) {
	m := protojson.MarshalOptions{
		UseProtoNames:   false,
		EmitUnpopulated: true,
	}

	w.Header().Set("Content-Type", "application/json")

	data, err := m.Marshal(&interestv1.ListInterestResponse{
		Interest: interests,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func getInterest(w http.ResponseWriter, r *http.Request) {
	m := protojson.MarshalOptions{
		UseProtoNames:   false,
		EmitUnpopulated: true,
	}

	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")

	for _, interest := range interests {
		if interest.Id == id {
			data, err := m.Marshal(interest)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(data)
			return
		}
	}

	http.Error(w, "interest not found", http.StatusNotFound)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/interest", listInterest)
	mux.HandleFunc("/interest/{id}", getInterest)

	log.Println("interest service listening on :8082")

	err := http.ListenAndServe(":8082", common.CorsMiddleware(mux))
	if err != nil {
		log.Fatal(err)
	}
}
