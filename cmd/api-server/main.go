package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tahanasir/dicom-service/internal/transport"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	router.Post("/v1/upload", transport.Upload())
	router.Get("/v1/extract", transport.Extract())
	router.Get("/v1/convert", transport.Convert())

	log.Fatal(http.ListenAndServe(":8080", router))
}
