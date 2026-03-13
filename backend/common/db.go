package common

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// OpenDB opens a database connection with a sensible connection pool configuration.
func OpenDB(driverName, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}

// QueryDB executes a query with an OTel span. The tracer name and peer.service are
// derived from OTEL_SERVICE_NAME (peer.service = "{OTEL_SERVICE_NAME}-db").
// The caller must call span.End() when done scanning rows.
func QueryDB(ctx context.Context, db *sql.DB, query string, args ...any) (*sql.Rows, trace.Span, error) {
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "unknown"
	}
	peerService := fmt.Sprintf("%s-db", serviceName)

	ctx, span := otel.Tracer(serviceName).Start(ctx, "db.query", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		attribute.String("db.statement", query),
		attribute.String("db.system", "mysql"),
		attribute.String("peer.service", peerService),
	)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.End()
		return nil, nil, err
	}

	return rows, span, nil
}
