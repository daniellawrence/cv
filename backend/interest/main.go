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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
		ctx, span := otel.Tracer("interest").Start(r.Context(), "db.query.interest", trace.WithSpanKind(trace.SpanKindClient))
		defer span.End()
		span.SetAttributes(
			attribute.String("db.statement", INTEREST_SQL),
			attribute.String("db.system", "mysql"),
			attribute.String("peer.service", "interest-db"),
		)

		rows, err := db.QueryContext(ctx, INTEREST_SQL)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

		ctx, span := otel.Tracer("interest").Start(r.Context(), "db.query.interest.byid", trace.WithSpanKind(trace.SpanKindClient))
		defer span.End()
		span.SetAttributes(
			attribute.String("db.statement", INTEREST_BY_ID_SQL),
			attribute.String("db.system", "mysql"),
			attribute.String("peer.service", "interest-db"),
		)

		rows, err := db.QueryContext(ctx, INTEREST_BY_ID_SQL, id)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
	fmt.Print("%s", DB_URL)
	db, err := sql.Open("mysql", DB_URL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer func() { _ = db.Close() }()

	mux := http.NewServeMux()

	mux.HandleFunc("/interest", listInterest(db))
	mux.HandleFunc("/interest/{id}", getInterest(db))
	log.Fatal(common.Listen(mux))
}
