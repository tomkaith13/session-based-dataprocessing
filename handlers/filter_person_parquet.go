package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/tomkaith13/session-based-dataprocessing/models"
)

func FilterPersonsParquetHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("duckdb", "")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`SELECT * FROM 'file.parquet' WHERE age < 90 AND age >= 50`)
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
