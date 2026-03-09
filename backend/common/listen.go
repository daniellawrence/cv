package common

import (
	"log"
	"net/http"
)

func Listen(mux *http.ServeMux) error {
	addr := GetListenAddr()
	log.Printf("Starting server on %s\n", addr)

	err := http.ListenAndServe(addr, CorsMiddleware(mux))
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
