package handlers

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"

	"github.com/tomkaith13/session-based-dataprocessing/models"
	"github.com/tomkaith13/session-based-dataprocessing/mongo"
)

func CreatePersonHandler(w http.ResponseWriter, r *http.Request) {
	mongoURI := os.Getenv("MONGO_URI")
	client, err := mongo.GetMongoClient(mongoURI)
	if err != nil {
		w.Write([]byte("Unable to connect to mongo" + err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	collection := client.Database("mydatabase").Collection("persons")

	for i := 0; i < 5; i++ {
		id := uuid.NewString()
		person := models.Person{
			Id:   id,
			Name: "name" + strconv.Itoa(i),
			Age:  21,
			City: "Toronto",
		}

		_, err := collection.InsertOne(context.Background(), person)
		if err != nil {
			w.Write([]byte("Unable to insert document:" + err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.Write([]byte("Inserted successfully"))
	w.WriteHeader(http.StatusOK)
}
