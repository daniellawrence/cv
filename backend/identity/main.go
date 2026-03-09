package main

import (
	"log"
	"net/http"

	"github.com/daniellawrence/cv/backend/common"
	identityv1 "github.com/daniellawrence/cv/gen/go/identity/v1"

	"google.golang.org/protobuf/encoding/protojson"
)

var identitys = []*identityv1.Identity{
	{
		Id:       "dsl",
		Name:     "Daniel Lawrence",
		Email:    "dslawrence@protonmail.com",
		Linkedin: "https://www.linkedin.com/in/dannylawrence/",
	},
}

func listidentity(w http.ResponseWriter, r *http.Request) {
	m := protojson.MarshalOptions{
		UseProtoNames:   false,
		EmitUnpopulated: true,
	}

	w.Header().Set("Content-Type", "application/json")

	data, err := m.Marshal(&identityv1.ListIdentityResponse{
		Identity: identitys,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func getIdentity(w http.ResponseWriter, r *http.Request) {
	m := protojson.MarshalOptions{
		UseProtoNames:   false,
		EmitUnpopulated: true,
	}

	w.Header().Set("Content-Type", "application/json")

	id := r.PathValue("id")

	for _, identity := range identitys {
		if identity.Id == id {
			data, err := m.Marshal(identity)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(data)
			return
		}
	}

	http.Error(w, "identity not found", http.StatusNotFound)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/identity", listidentity)
	mux.HandleFunc("/identity/{id}", getIdentity)

	log.Println("identity service listening on :8083")

	err := http.ListenAndServe(":8083", common.CorsMiddleware(mux))
	if err != nil {
		log.Fatal(err)
	}
}
