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

const (
	DefaultDBDriver = "mysql"
	DBURLEnvVar     = "DATABASE_URL"
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

	// Verify the connection works
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

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

// GetDatabaseURL retrieves the database URL from environment variable or returns a default fallback.
func GetDatabaseURL(defaultURL string) string {
	if url := os.Getenv(DBURLEnvVar); url != "" {
		return url
	}
	return defaultURL
}

// ValidateConnectionString checks if a connection string is valid by attempting to parse it.
// Returns true and no error if the connection string appears valid, false with an error otherwise.
func ValidateConnectionString(dsn string) (bool, error) {
	if dsn == "" {
		return false, fmt.Errorf("empty database connection string")
	}

	// Basic validation: check for common MySQL DSN patterns
	// Expected format: user@tcp(host:port)/database or similar
	parts := splitDSN(dsn)
	
	if len(parts) < 2 {
		return false, fmt.Errorf("invalid database connection string format: %s", dsn)
	}

	return true, nil
}

// ConnectWithValidation opens a database connection after validating the connection string.
// It uses the defaultURL as a fallback if no DATABASE_URL environment variable is set.
func ConnectWithValidation(defaultURL string) (*sql.DB, error) {
	dbURL := GetDatabaseURL(defaultURL)

	valid, err := ValidateConnectionString(dbURL)
	if !valid || err != nil {
		return nil, fmt.Errorf("invalid database connection string: %w", err)
	}

	db, err := OpenDB(DefaultDBDriver, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database with URL '%s': %w", dbURL, err)
	}

	return db, nil
}

// splitDSN splits a DSN string by common delimiters to validate its structure.
func splitDSN(dsn string) []string {
	var parts []string
	current := ""
	inQuotes := false
	afterAt := true // Track if we're after the @ symbol
	
	for _, ch := range dsn {
		if ch == '\'' || ch == '"' {
			inQuotes = !inQuotes
			continue
		}
		
		// After @, we're in host:port/database section
		if ch == '@' && !inQuotes {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
			afterAt = true
			continue
		}
		
		// After @, split on / to separate host:port from database name
		if ch == '/' && !inQuotes && afterAt {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
			continue
		}
		
		// After @, split on ? for query parameters
		if ch == '?' && !inQuotes && afterAt {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
			continue
		}
		
		// Keep track of colon in host:port (after @ and before /)
		if ch == ':' && !inQuotes && afterAt {
			current += string(ch)
			continue
		}
		
		current += string(ch)
	}
	
	if current != "" {
		parts = append(parts, current)
	}
	
	return parts
}
