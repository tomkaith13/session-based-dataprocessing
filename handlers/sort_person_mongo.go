package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/tomkaith13/session-based-dataprocessing/models"
	"github.com/tomkaith13/session-based-dataprocessing/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SortPersonsHandler(w http.ResponseWriter, r *http.Request) {
	mongoURI := os.Getenv("MONGO_URI")
	client, err := mongo.GetMongoClient(mongoURI)
	if err != nil {
		w.Write([]byte("Unable to connect to mongo" + err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	collection := client.Database("mydatabase").Collection("persons")

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"createdAt", 1}})
	findOptions.SetLimit(200)

	// sort people
	cursor, err := collection.Find(context.TODO(), bson.D{}, findOptions)
	// cursor, err := collection.Find(context.TODO(), bson.M{"city": "Toronto"})
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	// Decode results into a slice of Person structs
	var people []models.Person
	if err = cursor.All(context.TODO(), &people); err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// fmt.Println("people:", people)

	// Encode results as JSON and write to the response
	json.NewEncoder(w).Encode(people)
}
