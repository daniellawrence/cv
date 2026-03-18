package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/daniellawrence/cv/backend/common"
	identityv1 "github.com/daniellawrence/cv/gen/go/identity/v1"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	IDENTITY_SQL        string = "SELECT id, name, email, linkedin FROM identity"
	IDENTITY_BY_ID_SQL string = "SELECT id, name, email, linkedin FROM identity WHERE id = ?"
)

func listidentity(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, span, err := common.QueryDB(r.Context(), db, IDENTITY_SQL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer span.End()
		defer func() { _ = rows.Close() }()

		var results []*identityv1.Identity
		for rows.Next() {
			i := &identityv1.Identity{}
			if err := rows.Scan(&i.Id, &i.Name, &i.Email, &i.Linkedin); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			results = append(results, i)
		}

		m := protojson.MarshalOptions{
			UseProtoNames:   false,
			EmitUnpopulated: true,
		}

		w.Header().Set("Content-Type", "application/json")

		data, err := m.Marshal(&identityv1.ListIdentityResponse{
			Identity: results,
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
}

func getIdentity(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		rows, span, err := common.QueryDB(r.Context(), db, IDENTITY_BY_ID_SQL, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer span.End()
		defer func() { _ = rows.Close() }()

		m := protojson.MarshalOptions{
			UseProtoNames:   false,
			EmitUnpopulated: true,
		}

		w.Header().Set("Content-Type", "application/json")

		if rows.Next() {
			i := &identityv1.Identity{}
			if err := rows.Scan(&i.Id, &i.Name, &i.Email, &i.Linkedin); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data, err := m.Marshal(i)
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

		http.Error(w, "identity not found", http.StatusNotFound)
	}
}

func main() {
	defaultDBURL := "root@tcp(identity-db.identity:3306)/cv"
	
	db, err := common.ConnectWithValidation(defaultDBURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() { _ = db.Close() }()

	mux := http.NewServeMux()

	mux.HandleFunc("/identity", listidentity(db))
	mux.HandleFunc("/identity/{id}", getIdentity(db))

	log.Fatal(common.Listen(mux))
}
