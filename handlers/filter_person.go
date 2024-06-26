package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/tomkaith13/session-based-dataprocessing/models"
	"github.com/tomkaith13/session-based-dataprocessing/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FilterPersonsHandler(w http.ResponseWriter, r *http.Request) {
	mongoURI := os.Getenv("MONGO_URI")
	client, err := mongo.GetMongoClient(mongoURI)
	if err != nil {
		w.Write([]byte("Unable to connect to mongo" + err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	collection := client.Database("mydatabase").Collection("persons")

	lowerInc := rand.Intn(100000)
	lowerInc += 25000

	higherInc := rand.Intn(500000)
	higherInc += 25000

	filters := []models.Filter{
		// {Field: "user_id", Ops: []models.Operator{{Op: "", Value: "3"}}},
		{Field: "age", Ops: []models.Operator{{Op: "$gte", Value: 20}, {Op: "$lte", Value: 70}}},
		{Field: "inc", Ops: []models.Operator{{Op: "$gte", Value: lowerInc}, {Op: "$lte", Value: higherInc}}},
		// {Field: "city", Ops: []models.Operator{{Op: "$in", Value: []string{"Toronto", "Bangalore", "Mountain View"}}}},
		// {Field: "user_id", Ops: []models.Operator{{Op: "$in", Value: []string{"3", "5", "100", "50000", "-1"}}}}, // We add various valid and invalid ids
	}

	filter := buildFilter(filters)
	fmt.Println("filter:", filter)

	findOptions := options.Find()
	findOptions.SetLimit(100)

	// Find matching people in the collection
	cursor, err := collection.Find(context.TODO(), filter, findOptions)
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

// buildFilter constructs a bson.M filter based on the provided Filter conditions
func buildFilter(filters []models.Filter) bson.M {
	bsonFilter := bson.M{}

	for _, f := range filters {
		if len(f.Ops) == 1 {
			onlyOp := f.Ops[0]
			if onlyOp.Op == "" {
				// If no operator, assume equality
				bsonFilter[f.Field] = onlyOp.Value
			} else {
				// Use the specified operator
				bsonFilter[f.Field] = bson.M{onlyOp.Op: onlyOp.Value}
			}
			continue
		}
		// otherwise we use an AND
		opMap := bson.M{}
		for _, op := range f.Ops {
			opMap[op.Op] = op.Value
		}
		bsonFilter[f.Field] = opMap

	}
	return bsonFilter
}
