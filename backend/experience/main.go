package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	common "github.com/daniellawrence/cv/backend/common"
	experiencev1 "github.com/daniellawrence/cv/gen/go/experience/v1"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	EXPERIENCE_SQL string = "SELECT id, company, title, start_date, end_date, location, highlights, skills FROM experience"
)

func listExperience(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		offset := common.GetPathInt(r, "offset", 0)
		limit := common.GetPathInt(r, "limit", 4)

		rows, span, err := common.QueryDB(r.Context(), db, EXPERIENCE_SQL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer span.End()
		defer func() { _ = rows.Close() }()

		var results []*experiencev1.Experience
		for rows.Next() {
			e := &experiencev1.Experience{}
			var highlightsJSON, skillsJSON string
			if err := rows.Scan(&e.Id, &e.Company, &e.Title, &e.StartDate, &e.EndDate, &e.Location, &highlightsJSON, &skillsJSON); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := json.Unmarshal([]byte(highlightsJSON), &e.Highlights); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := json.Unmarshal([]byte(skillsJSON), &e.Skills); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			results = append(results, e)
		}

		page := results
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
}

func main() {
	defaultDBURL := "root@tcp(experience-db.experience:3306)/cv"
	
	db, err := common.ConnectWithValidation(defaultDBURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() { _ = db.Close() }()

	mux := http.NewServeMux()

	mux.HandleFunc("/experience", listExperience(db))
	mux.HandleFunc("/experience/{offset}/{limit}", listExperience(db))
	log.Fatal(common.Listen(mux))
}
