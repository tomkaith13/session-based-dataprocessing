package models

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	Id        string    `bson:"_id" json:"id"`
	UserId    string    `bson:"user_id" json:"user-id"`
	Name      string    `bson:"name" json:"Name"`
	City      string    `bson:"city" json:"City"`
	Age       int       `bson:"age" json:"Age"`
	CreatedAt time.Time `bson:"createdAt" json:"Created At"`
}

// Filter struct to represent a single filter condition
type Filter struct {
	Field string
	Ops   []Operator
}
type Operator struct {
	Op    string // Comparison operator: "$eq", "$ne", "$gt", "$gte", "$lt", "$lte", etc.
	Value interface{}
}

type PersonParquet struct {
	Id       uuid.UUID `parquet:"id"`
	Name     string    `parquet:"name,lz4"`
	Age      int       `parquet:"age"`
	Location string    `parquet:"location,lz4"`
}
