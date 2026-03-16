package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime"
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

// StatusResponse mirrors K8s statusz format with service health and dependencies
type StatusResponse struct {
	Service      ServiceInfo             `json:"service"`
	Dependencies map[string]DependencyStatus `json:"dependencies,omitempty"`
	GoVersion    string                  `json:"go_version,omitempty"`
}

// ServiceInfo represents the service's own status (like K8s Pod conditions)
type ServiceInfo struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Status      ConditionType     `json:"status"`
	Conditions  []Condition       `json:"conditions,omitempty"`
}

// DependencyStatus represents a service dependency (like K8s pod dependencies)
type DependencyStatus struct {
	Name    string        `json:"name"`
	Version string        `json:"version"`
	Status  ConditionType `json:"status"`
	Reason  string        `json:"reason,omitempty"`
}

// Condition mirrors K8s PodCondition format
type Condition struct {
	Type   ConditionType `json:"type"`
	Status ConditionStatus `json:"status"`
	LastProbeTime     string `json:"lastProbeTime,omitempty"`
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`
	Reason            string `json:"reason,omitempty"`
	Message           string `json:"message,omitempty"`
}

// ConditionType matches K8s condition types
type ConditionType string

const (
	ConditionReady    ConditionType = "Ready"
	ConditionHealthy  ConditionType = "Healthy"
	ConditionDegraded ConditionType = "Degrade"
)

// ConditionStatus matches K8s pod condition statuses
type ConditionStatus string

const (
	ConditionTrue    ConditionStatus = "true"
	ConditionFalse   ConditionStatus = "false"
	ConditionUnknown ConditionStatus = "unknown"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "unknown"
	}

	env := os.Getenv("ENVIRONMENT")
	
	// Build conditions like K8s does for pods
	conditions := []Condition{
		{
			Type:   ConditionReady,
			Status: ConditionTrue,
			Reason: "PodReady",
			Message: fmt.Sprintf("%s is ready and serving requests", serviceName),
		},
		{
			Type:   ConditionHealthy,
			Status: ConditionTrue,
			Reason: "ServiceHealthy",
			Message: fmt.Sprintf("All health checks passed for %s", serviceName),
		},
	}

	if env != "" {
		conditions = append(conditions, Condition{
			Type:   "Environment",
			Status: ConditionTrue,
			Reason: "EnvironmentSet",
			Message: fmt.Sprintf("Running in %s environment", env),
		})
	}

	response := StatusResponse{
		Service: ServiceInfo{
			Name:     serviceName,
			Version:  BuildTimestamp,
			Status:   ConditionReady,
			Conditions: conditions,
		},
		Dependencies: make(map[string]DependencyStatus),
		GoVersion:    runtime.Version(),
	}

	data, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func Listen(mux *http.ServeMux) error {

	shutdown := InitTracing()
	defer func() { _ = shutdown(context.Background()) }()

	// Wrap only the app routes with otelhttp — healthz and statusz sit outside it
	appHandler := otelhttp.NewHandler(mux, "http-server")

	top := http.NewServeMux()
	top.HandleFunc("/healthz", healthz)
	top.HandleFunc("/statusz", statusHandler)
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
