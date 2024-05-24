package handlers

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	mongo_drv "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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

	//  As per docs, CreateIndex is not applied if one already is in-place. See https://stackoverflow.com/a/49476360/224640
	indexModel := mongo_drv.IndexModel{
		Keys:    bson.M{"createdAt": 1},
		Options: options.Index().SetExpireAfterSeconds(600), // Expire after 10 min
	}
	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userIdIndexModel := mongo_drv.IndexModel{
		Keys: bson.M{"user_id": 1},
	}
	_, err = collection.Indexes().CreateOne(context.Background(), userIdIndexModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ageIndexModel := mongo_drv.IndexModel{
		Keys: bson.M{"age": 1},
	}
	_, err = collection.Indexes().CreateOne(context.Background(), ageIndexModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currTime := time.Now()
	people := []any{}
	// Now assume we get this from a bulk-load operation either from a network call or a blob storage read.
	for i := 1; i <= 1000000; i++ {
		id := uuid.NewString()

		randAge := rand.Intn(91)
		randAge += 10
		person := models.Person{
			Id:        id,
			UserId:    strconv.Itoa(i),
			Name:      "name" + strconv.Itoa(i),
			Age:       randAge,
			City:      "Toronto",
			CreatedAt: currTime,
		}
		people = append(people, person)
	}
	_, err = collection.InsertMany(context.Background(), people)
	if err != nil {
		w.Write([]byte("error inserting many:" + err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Inserted successfully"))
	w.WriteHeader(http.StatusOK)
}
