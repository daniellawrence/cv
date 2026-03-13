package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/daniellawrence/cv/backend/common"
	identityv1 "github.com/daniellawrence/cv/gen/go/identity/v1"

	_ "github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	IDENTITY_SQL       string = "SELECT id, name, email, linkedin FROM identity"
	IDENTITY_BY_ID_SQL string = "SELECT id, name, email, linkedin FROM identity WHERE id = ?"
)

func listidentity(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer("identity").Start(r.Context(), "db.query.identity", trace.WithSpanKind(trace.SpanKindClient))
		defer span.End()
		span.SetAttributes(
			attribute.String("db.statement", IDENTITY_SQL),
			attribute.String("db.system", "mysql"),
			attribute.String("peer.service", "identity-db"),
		)

		rows, err := db.QueryContext(ctx, IDENTITY_SQL)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

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

		ctx, span := otel.Tracer("identity").Start(r.Context(), "db.query.identity.byid", trace.WithSpanKind(trace.SpanKindClient))
		defer span.End()
		span.SetAttributes(
			attribute.String("db.statement", IDENTITY_BY_ID_SQL),
			attribute.String("db.system", "mysql"),
			attribute.String("peer.service", "identity-db"),
		)

		rows, err := db.QueryContext(ctx, IDENTITY_BY_ID_SQL, id)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

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
	DB_URL := "root@tcp(identity-db.identity:3306)/cv"
	fmt.Printf(DB_URL)
	db, err := sql.Open("mysql", DB_URL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/identity", listidentity(db))
	mux.HandleFunc("/identity/{id}", getIdentity(db))

	common.Listen(mux)
}
