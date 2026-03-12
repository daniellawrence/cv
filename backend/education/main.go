package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	common "github.com/daniellawrence/cv/backend/common"
	educationv1 "github.com/daniellawrence/cv/gen/go/education/v1"

	_ "github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	EDUCATION_SQL string = "SELECT id, institution, degree, field_of_study, start_date, end_date FROM education"
)

func listEducation(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer("education").Start(r.Context(), "db.query.education", trace.WithSpanKind(trace.SpanKindClient))
		defer span.End()
		span.SetAttributes(
			attribute.String("db.statement", EDUCATION_SQL),
			attribute.String("db.system", "mysql"),
			attribute.String("peer.service", "education-db"),
		)

		rows, err := db.QueryContext(ctx, EDUCATION_SQL)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

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
	DB_URL := "root@tcp(education-db.education:3306)/cv"
	fmt.Printf(DB_URL)
	db, err := sql.Open("mysql", DB_URL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/education", listEducation(db))

	common.Listen(mux)
}
