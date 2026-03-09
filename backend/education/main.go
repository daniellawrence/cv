package main

import (
	"net/http"

	common "github.com/daniellawrence/cv/backend/common"
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

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/education", listEducation)

	common.Listen(mux)
}
