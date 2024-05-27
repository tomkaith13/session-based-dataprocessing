package handlers

import (
	"context"
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
	// SELECT * FROM read_parquet('https://some.url/some_file.parquet');
	// See https://duckdb.org/docs/data/parquet/overview.html#examples
	rows, err := conn.QueryContext(ctx, `SELECT * FROM 'file.parquet' WHERE age < 90 AND age >= 50`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := []models.PersonParquet{}

	for rows.Next() {
		pperson := new(models.PersonParquet)
		err := rows.Scan(&pperson.Id, &pperson.Name, &pperson.Age, &pperson.Location)
		if err != nil {
			fmt.Println("unable to scan row", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		res = append(res, *pperson)
	}

	fmt.Println("count: ", len(res))

	w.WriteHeader(http.StatusAccepted)
}