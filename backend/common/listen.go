package common

import (
	"log"
	"net/http"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	_, err := w.Write([]byte("ok"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func zPages(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", healthz)
}

func Listen(mux *http.ServeMux) error {

	zPages(mux)

	addr := GetListenAddr()
	log.Printf("Starting server on %s\n", addr)

	err := http.ListenAndServe(addr, CorsMiddleware(mux))
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
