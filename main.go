package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tomkaith13/session-based-dataprocessing/handlers"
)

func main() {
	r := chi.NewRouter()

	// Middleware for logging and recovery
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Dummy Endpoint
	r.Get("/dummy", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a dummy endpoint!!!!"))
	})

	r.Post("/person", handlers.CreatePersonHandler)
	r.Post("/person-search", handlers.FilterPersonsHandler)
	r.Post("/person-parquet", handlers.CreatePersonParquetHandler)
	r.Post("/person-search-parquet", handlers.FilterPersonsParquetHandler)

	// Start Server
	http.ListenAndServe(":8080", r)
}
