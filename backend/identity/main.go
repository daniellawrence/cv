package main

import (
	"net/http"

	"github.com/daniellawrence/cv/backend/common"
	identityv1 "github.com/daniellawrence/cv/gen/go/identity/v1"

	"google.golang.org/protobuf/encoding/protojson"
)

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

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

			_, err = w.Write(data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}

	http.Error(w, "identity not found", http.StatusNotFound)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/identity", listidentity)
	mux.HandleFunc("/identity/{id}", getIdentity)

	common.Listen(mux)
}
