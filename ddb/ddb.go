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
}

func (ddb DuckDBService) GetDB() *sql.DB {
	return ddb.DB
}
