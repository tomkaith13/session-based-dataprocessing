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
	// Install httpfs extension - this is needed for loading remote files. See https://duckdb.org/docs/extensions/httpfs/overview.html
	if _, err := duckdb.Exec("INSTALL httpfs; LOAD httpfs;"); err != nil {
		log.Fatal(err)
	}
}

func (ddb DuckDBService) GetDB() *sql.DB {
	return ddb.DB
}
