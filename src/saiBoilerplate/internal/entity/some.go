package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Some struct {
	ID  primitive.ObjectID `bson:"_id"`
	Key string             `bson:"key"`
}
