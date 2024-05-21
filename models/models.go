package models

import "time"

type Person struct {
	Id        string    `bson:"_id"`
	Name      string    `bson:"name"`
	City      string    `bson:"city"`
	Age       int       `bson:"age"`
	CreatedAt time.Time `bson:"createdAt"`
}
