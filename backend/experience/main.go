package main

import (
	"net/http"

	"github.com/daniellawrence/cv/backend/common"
	experiencev1 "github.com/daniellawrence/cv/gen/go/experience/v1"

	"google.golang.org/protobuf/encoding/protojson"
)

func listExperience(w http.ResponseWriter, r *http.Request) {
	offset := common.GetPathInt(r, "offset", 0)
	limit := common.GetPathInt(r, "limit", 4)

	page := experiences
	if offset < len(page) {
		page = page[offset:]
	} else {
		page = nil
	}
	if limit < len(page) {
		page = page[:limit]
	}

	m := protojson.MarshalOptions{
		UseProtoNames:   false,
		EmitUnpopulated: true,
	}

	w.Header().Set("Content-Type", "application/json")

	data, err := m.Marshal(&experiencev1.ListExperienceResponse{
		Experience: page,
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
	mux.HandleFunc("/experience/{offset}/{limit}", listExperience)
	common.Listen(mux)
}
