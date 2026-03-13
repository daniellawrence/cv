package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	common "github.com/daniellawrence/cv/backend/common"
	interestv1 "github.com/daniellawrence/cv/gen/go/interest/v1"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	INTEREST_SQL       string = "SELECT id, type, names FROM interest"
	INTEREST_BY_ID_SQL string = "SELECT id, type, names FROM interest WHERE id = ?"
)

func scanInterest(rows *sql.Rows) (*interestv1.Interest, error) {
	i := &interestv1.Interest{}
	var namesJSON string
	if err := rows.Scan(&i.Id, &i.Type, &namesJSON); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(namesJSON), &i.Names); err != nil {
		return nil, err
	}
	return i, nil
}

func listInterest(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, span, err := common.QueryDB(r.Context(), db, INTEREST_SQL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer span.End()
		defer func() { _ = rows.Close() }()

		var results []*interestv1.Interest
		for rows.Next() {
			i, err := scanInterest(rows)
			if err != nil {
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

		data, err := m.Marshal(&interestv1.ListInterestResponse{
			Interest: results,
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

func getInterest(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		rows, span, err := common.QueryDB(r.Context(), db, INTEREST_BY_ID_SQL, id)
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
			i, err := scanInterest(rows)
			if err != nil {
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

		http.Error(w, "interest not found", http.StatusNotFound)
	}
}

func main() {
	DB_URL := "root@tcp(interest-db.interest:3306)/cv"
	fmt.Printf("%s\n", DB_URL)
	db, err := common.OpenDB("mysql", DB_URL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() { _ = db.Close() }()

	mux := http.NewServeMux()

	mux.HandleFunc("/interest", listInterest(db))
	mux.HandleFunc("/interest/{id}", getInterest(db))
	log.Fatal(common.Listen(mux))
}
