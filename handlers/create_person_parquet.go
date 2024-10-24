package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/parquet-go/parquet-go"
	"github.com/tomkaith13/session-based-dataprocessing/models"
	"github.com/tomkaith13/session-based-dataprocessing/utils"
)

const (
	parquetFilePath = "./file.parquet" // this filename can be renamed with the sessionid to ensure we have one file per session
	TTL             = time.Duration(15) * time.Minute
)

func CreatePersonParquetHandler(w http.ResponseWriter, r *http.Request) {

	f, err := os.Create(parquetFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	persons := []models.PersonParquet{}
	writer := parquet.NewGenericWriter[models.PersonParquet](f)

	buffer := parquet.NewGenericBuffer[models.PersonParquet](
		parquet.SortingRowGroupConfig(
			parquet.SortingColumns(
				parquet.Ascending("createdAt"),
			),
		),
	)

	for i := 1; i < 5000000; i++ {
		id := uuid.New()

		randAge := rand.Intn(91)
		randAge += 10

		// Income is randomly generated between 25k and 1.5M.
		// This is to accomodate high cardinality columns
		randIncome := rand.Intn(1500000 - 25000)
		randIncome += 25000

		person := models.PersonParquet{
			Id:       id,
			UserId:   strconv.Itoa(i),
			Name:     "name" + strconv.Itoa(i),
			Age:      randAge,
			Location: utils.RandomizedLocationCreator(),
			Created:  time.Now().Add(time.Duration(randIncome*24) * time.Hour),
			Income:   randIncome,
		}
		persons = append(persons, person)

	}
	// _, err = writer.Write(persons)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	sort.Sort(buffer)
	buffer.Write(persons)
	_, err = parquet.CopyRows(writer, buffer.Rows())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = writer.Close()
	_ = f.Close()

	rf, _ := os.Open(f.Name())
	fi, _ := rf.Stat()

	fmt.Println("size:", fi.Size())
	fmt.Println("path:", fi.Name())

	// Now we can send this to a blob storage like GCS with an Object Lifecycle to enforce longer TTL if we want or use BQ to filter from these directly

	// For now, lets assume this happens on disk and goruntime is the one who cleansup after
	time.AfterFunc(TTL, func() {

		err := os.Remove(parquetFilePath)
		if err != nil {
			fmt.Println("error removed parquet file")
		} else {
			fmt.Println("removed parquet file")
		}
	})

	w.WriteHeader(http.StatusCreated)

}
