package main

import (
	"database/sql"
	"log"
	"net/http"

	common "github.com/daniellawrence/cv/backend/common"
	educationv1 "github.com/daniellawrence/cv/gen/go/education/v1"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	EDUCATION_SQL string = "SELECT id, institution, degree, field_of_study, start_date, end_date FROM education"
)

func listEducation(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, span, err := common.QueryDB(r.Context(), db, EDUCATION_SQL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer span.End()
		defer func() { _ = rows.Close() }()

		var results []*educationv1.Education
		for rows.Next() {
			e := &educationv1.Education{}
			if err := rows.Scan(&e.Id, &e.Institution, &e.Degree, &e.FieldOfStudy, &e.StartDate, &e.EndDate); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			results = append(results, e)
		}

		m := protojson.MarshalOptions{
			UseProtoNames:   false,
			EmitUnpopulated: true,
		}

		w.Header().Set("Content-Type", "application/json")

		data, err := m.Marshal(&educationv1.ListEducationResponse{
			Education: results,
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
	defaultDBURL := "root@tcp(education-db.education:3306)/cv"
	
	db, err := common.ConnectWithValidation(defaultDBURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() { _ = db.Close() }()

	mux := http.NewServeMux()

	mux.HandleFunc("/education", listEducation(db))

	log.Fatal(common.Listen(mux))
}
