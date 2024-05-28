package ddb

import (
	"database/sql"
	"log"
)

type DuckDBService struct {
	DB *sql.DB
}

var DService DuckDBService

func InitDuckDB() {
	duckdb, err := sql.Open("duckdb", "")
	if err != nil {
		log.Fatal(err)
		return

	}

	DService.DB = duckdb
	// Install httpfs extension
	if _, err := duckdb.Exec("INSTALL httpfs; LOAD httpfs;"); err != nil {
		log.Fatal(err)
	}
}

func (ddb DuckDBService) GetDB() *sql.DB {
	return ddb.DB
}
