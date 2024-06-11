package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tomkaith13/session-based-dataprocessing/ddb"
	"github.com/tomkaith13/session-based-dataprocessing/handlers"
)

func main() {
	r := chi.NewRouter()

	// Middleware for logging and recovery
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// initialize duckdb
	ddb.InitDuckDB()

	// Dummy Endpoint
	r.Get("/dummy", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a dummy endpoint!!!!"))
	})

	r.Post("/person", handlers.CreatePersonHandler)
	r.Post("/person-search", handlers.FilterPersonsHandler)
	r.Post("/person-parquet", handlers.CreatePersonParquetHandler)
	r.Post("/person-search-parquet", handlers.FilterPersonsParquetHandler)

	r.Post("/person-parquet-table", handlers.CreatepersonParquetTableHandler)
	r.Post("/person-search-parquet-table", handlers.FilterPersonsParquetTableHandler)

	// Start Server
	http.ListenAndServe(":8080", r)
}
