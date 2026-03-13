package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	common "github.com/daniellawrence/cv/backend/common"
	experiencev1 "github.com/daniellawrence/cv/gen/go/experience/v1"

	_ "github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	EXPERIENCE_SQL string = "SELECT id, company, title, start_date, end_date, location, highlights, skills FROM experience"
)

func listExperience(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		offset := common.GetPathInt(r, "offset", 0)
		limit := common.GetPathInt(r, "limit", 4)

		ctx, span := otel.Tracer("experience").Start(r.Context(), "db.query.experience", trace.WithSpanKind(trace.SpanKindClient))
		defer span.End()
		span.SetAttributes(
			attribute.String("db.statement", EXPERIENCE_SQL),
			attribute.String("db.system", "mysql"),
			attribute.String("peer.service", "experience-db"),
		)

		rows, err := db.QueryContext(ctx, EXPERIENCE_SQL)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
	DB_URL := "root@tcp(experience-db.experience:3306)/cv"
	fmt.Printf("%s\n", DB_URL)
	db, err := sql.Open("mysql", DB_URL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() { _ = db.Close() }()

	mux := http.NewServeMux()

	mux.HandleFunc("/experience", listExperience(db))
	mux.HandleFunc("/experience/{offset}/{limit}", listExperience(db))
	log.Fatal(common.Listen(mux))
}
