package main

import (
	"net/http"

	common "github.com/daniellawrence/cv/backend/common"
	interestv1 "github.com/daniellawrence/cv/gen/go/interest/v1"

	"google.golang.org/protobuf/encoding/protojson"
)

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

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

			_, err = w.Write(data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}

	http.Error(w, "interest not found", http.StatusNotFound)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/interest", listInterest)
	mux.HandleFunc("/interest/{id}", getInterest)
	common.Listen(mux)

}
