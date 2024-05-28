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

	// In theory, if we ever want this from GCS, we can do something like this
	// Install httpfs and configure credentials
	// if _, err := db.Exec(`
	//     INSTALL httpfs;
	//     LOAD httpfs;
	//     CREATE SECRET "gcs_credentials" (
	//         TYPE GCS,
	//         KEY_ID '<your-key-id>',
	//         SECRET '<your-secret-key>'
	//     );
	// `); err != nil {
	//     log.Fatal(err)
	// }
}

func (ddb DuckDBService) GetDB() *sql.DB {
	return ddb.DB
}
