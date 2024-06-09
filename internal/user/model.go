package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)



type User struct {
	ID      	primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username	string             `json:"username" bson:"username"`
	CreatedAt  primitive.DateTime `json:"created_at" bson:"created_at"`
}