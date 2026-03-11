package common

import (
	"context"
	"log"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

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
	defer shutdown(context.Background())

	// Wrap only the app routes with otelhttp — healthz sits outside it
	appHandler := otelhttp.NewHandler(mux, "http-server")

	top := http.NewServeMux()
	top.HandleFunc("/healthz", healthz)
	top.Handle("/", appHandler)

	addr := GetListenAddr()
	log.Printf("Starting server on %s\n", addr)

	err := http.ListenAndServe(addr, CorsMiddleware(top))
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
