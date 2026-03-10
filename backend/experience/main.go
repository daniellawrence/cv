package main

import (
	"net/http"

	"github.com/daniellawrence/cv/backend/common"
	experiencev1 "github.com/daniellawrence/cv/gen/go/experience/v1"

	"google.golang.org/protobuf/encoding/protojson"
)

func listExperience(w http.ResponseWriter, r *http.Request) {
	m := protojson.MarshalOptions{
		UseProtoNames:   false,
		EmitUnpopulated: true,
	}

	w.Header().Set("Content-Type", "application/json")

	data, err := m.Marshal(&experiencev1.ListExperienceResponse{
		Experience: experiences,
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

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/experience", listExperience)
	common.Listen(mux)
}
