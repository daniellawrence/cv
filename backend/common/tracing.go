package common

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
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

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp.Shutdown
}
