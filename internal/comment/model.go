package comment

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID      	primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TargetID 	primitive.ObjectID `json:"target_id" bson:"target_id"`
	AuthorID 	string `json:"author_id" bson:"author_id"`
	Content 	string `json:"content" bson:"content"`
	Replies		[]Comment `json:"replies" bson:"replies,omitempty"`
	CreatedAt  primitive.DateTime `json:"created_at" bson:"created_at"`
}

type CommentView struct {
	ID      	string `json:"id"`
	TargetID 	string `json:"target_id"`
	AuthorID 	string `json:"author_id"`
	Content 	string `json:"content"`
	Replies		[]CommentView `json:"replies"`
	CreatedAt  string `json:"created_at"`
}