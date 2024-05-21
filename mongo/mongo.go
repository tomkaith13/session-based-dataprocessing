package mongo

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoClient(mongoURI string) (*mongo.Client, error) {
	user := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	passwd := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	fmt.Println("user:", user)
	fmt.Println("passwd:", passwd)
	clientOptions := options.Client().ApplyURI(mongoURI).SetAuth(options.Credential{
		Username:   user,
		Password:   passwd,
		AuthSource: "admin",
	})
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return client, err
}
