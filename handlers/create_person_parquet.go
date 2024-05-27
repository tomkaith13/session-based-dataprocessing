package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/parquet-go/parquet-go"
)

type PersonParquet struct {
	Id       uuid.UUID `parquet:"id"`
	Name     string    `parquet:"name,lz4"`
	Age      int       `parquet:"age"`
	Location string    `parquet:"location,lz4"`
}

const (
	parquetFilePath = "./file.pq" // this filename can be renamed with the sessionid to ensure we have one file per session
	TTL             = time.Duration(1) * time.Minute
)

func CreatePersonParquetHandler(w http.ResponseWriter, r *http.Request) {

	f, _ := os.Create(parquetFilePath)
	writer := parquet.NewWriter(f)

	for i := 1; i < 1000000; i++ {
		id := uuid.New()

		randAge := rand.Intn(91)
		randAge += 10
		person := PersonParquet{
			Id:       id,
			Name:     "name" + strconv.Itoa(i),
			Age:      randAge,
			Location: "Toronto",
		}

		err := writer.Write(person)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	_ = writer.Close()
	_ = f.Close()

	rf, _ := os.Open(f.Name())
	fi, _ := rf.Stat()

	fmt.Println("size:", fi.Size())

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
