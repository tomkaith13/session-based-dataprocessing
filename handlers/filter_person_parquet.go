package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/tomkaith13/session-based-dataprocessing/ddb"
	"github.com/tomkaith13/session-based-dataprocessing/models"
)

func FilterPersonsParquetHandler(w http.ResponseWriter, r *http.Request) {

	db := ddb.DService.GetDB()
	err := db.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conn, err := db.Conn(context.Background())
	defer conn.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	//  We can also use read_parquet to read files directly from GCS and query from disk using something like:
	// rows, err := conn.QueryContext(ctx, `
	// SELECT name,age FROM 'https://github.com/tomkaith13/session-based-dataprocessing/raw/main/file.parquet'
	// WHERE age < 90 AND age >= 30 AND userId IN (1,10,100,500,1000,50000, 10000,100000, 500000)
	// `)
	// See https://duckdb.org/docs/data/parquet/overview.html#examples
	rows, err := conn.QueryContext(ctx, `
	SELECT name,age 
	FROM 'file.parquet' 
	WHERE age < 90 AND age >= 50 AND userId IN (1,10,100,500,1000,50000, 10000,100000, 500000)
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := []models.PersonParquet{}

	for rows.Next() {
		pperson := new(models.PersonParquet)
		// err := rows.Scan(&pperson.Id, &pperson.Name)
		// Also read up Partial Reads: https://duckdb.org/docs/data/parquet/overview.html#partial-reading
		err := rows.Scan(&pperson.Name, &pperson.Age)
		if err != nil {
			fmt.Println("unable to scan row", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		res = append(res, *pperson)
	}

	fmt.Println("count: ", len(res))
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(res)

	w.WriteHeader(http.StatusAccepted)
}
