package project

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)



type Project struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title    string             `json:"title" bson:"title"`
	Content  string             `json:"content" bson:"content"`
	AuthorID primitive.ObjectID `json:"author_id" bson:"author_id"`
	Link	 string             `json:"link" bson:"link"`
	Tags	 string           	`json:"tags" bson:"tags"`
	Votes	 int            	`json:"votes" bson:"votes"`
}