package common

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// Set via -ldflags at build time from Dockerfile build args.
var (
	ContentSHA1    = "unknown"
	BuildTimestamp = "unknown"
)

func InitTracing() func(context.Context) error {
	ctx := context.Background()

	// Automatically reads:
	// OTEL_EXPORTER_OTLP_ENDPOINT
	// OTEL_EXPORTER_OTLP_HEADERS
	// OTEL_EXPORTER_OTLP_PROTOCOL
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		panic(err)
	}

	res := resource.NewWithAttributes(
		resource.Default().SchemaURL(),
		attribute.String("image.content_sha1", ContentSHA1),
		attribute.String("image.build_timestamp", BuildTimestamp),
	)

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp.Shutdown
}
