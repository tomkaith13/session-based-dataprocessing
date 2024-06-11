package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/tomkaith13/session-based-dataprocessing/ddb"
)

func CreatePersonParquetViewHandler(w http.ResponseWriter, r *http.Request) {
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
	CREATE OR REPLACE VIEW personParquetView AS 
	SELECT * FROM 'https://github.com/tomkaith13/session-based-dataprocessing/raw/main/file.parquet'
	`
	_, err = conn.QueryContext(ctx, query)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}
