package common

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func accessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		spanCtx := trace.SpanFromContext(r.Context()).SpanContext()
		slog.InfoContext(r.Context(), "access",
			slog.String("http.request.method", r.Method),
			slog.String("url.path", r.URL.Path),
			slog.Int("http.response.status_code", rec.status),
			slog.Float64("http.server.request.duration", time.Since(start).Seconds()),
			slog.String("trace_id", spanCtx.TraceID().String()),
			slog.String("span_id", spanCtx.SpanID().String()),
		)
	})
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	_, err := w.Write([]byte("ok"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Listen(mux *http.ServeMux) error {

	shutdown := InitTracing()
	defer func() { _ = shutdown(context.Background()) }()

	// Wrap only the app routes with otelhttp — healthz sits outside it
	appHandler := otelhttp.NewHandler(mux, "http-server")

	top := http.NewServeMux()
	top.HandleFunc("/healthz", healthz)
	top.Handle("/", appHandler)

	addr := GetListenAddr()
	log.Printf("Starting server on %s content_sha1=%s build_timestamp=%s\n", addr, ContentSHA1, BuildTimestamp)

	err := http.ListenAndServe(addr, accessLog(CorsMiddleware(top)))
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
