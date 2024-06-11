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
	Income    int       `bson:"inc" json:"income"`
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
	Id       uuid.UUID `parquet:"id" json:"-"`
	UserId   string    `parquet:"userId" json:"-"`
	Name     string    `parquet:"name,snappy" json:"Name"`
	Age      int       `parquet:"age" json:"Age"`
	Location string    `parquet:"location,snappy" json:"-"`
	Income   int       `parquet:"inc,snappy" json:"income"`
	Created  time.Time `parquet:"createdAt" json:"created-at,omitempty"`
}
