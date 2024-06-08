package project

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title    string             `json:"title" bson:"title"`
	Content  string             `json:"content" bson:"content"`
	AuthorID string 			`json:"author_id" bson:"author_id,omitempty"`
	Link	 string             `json:"link" bson:"link"`
	Tags	 []string            `json:"tags" bson:"tags,omitempty"`
	Votes	 int32            	`json:"votes" bson:"votes,omitempty"`
	Voted	 bool            	`json:"voted" bson:"voted,omitempty"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
}

type ProjectView struct {
	ID      string `json:"id" bson:"_id,omitempty"`
	Title    string `json:"title" bson:"title"`
	Content  string `json:"content" bson:"content"`
	AuthorID string `json:"author_id" bson:"author_id,omitempty"`
	Link	 string `json:"link" bson:"link"`
	Tags	 []string `json:"tags" bson:"tags,omitempty"`
	Votes	 int32 `json:"votes" bson:"votes,omitempty"`
	Voted	 bool  `json:"voted" bson:"voted,omitempty"`
	CreatedAt string `json:"created_at" bson:"created_at"`
}

type ProjectVote struct {
	ProjectID primitive.ObjectID `json:"project_id" bson:"project_id"`
	AuthorID  string             `json:"author_id" bson:"author_id"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
}