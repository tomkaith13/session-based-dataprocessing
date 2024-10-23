package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/tomkaith13/session-based-dataprocessing/ddb"
)

func CreatepersonParquetTableHandler(w http.ResponseWriter, r *http.Request) {
	db := ddb.DService.GetDB()
	err := db.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conn, err := db.Conn(context.Background())
	if err != nil {
		err = errors.New("unable to connect to DB:" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	query := `
	CREATE OR REPLACE TABLE personParquetTable AS 
	SELECT * FROM 'file.parquet'
	`
	_, err = conn.QueryContext(ctx, query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	query2 := `
	CREATE INDEX createdAtIndex ON personParquetTable(createdAt)
	`
	_, err = conn.QueryContext(ctx, query2)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}
