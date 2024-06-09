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
	Comments int 				`json:"comments" bson:"comments,omitempty"`
	Awards  []Award 			`json:"awards" bson:"awards,omitempty"`
	AwardsTotal  int 			`json:"awards_total" bson:"awards_total,omitempty"`
}

type Award struct {
	ID      	primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TargetID 	primitive.ObjectID `json:"target_id" bson:"target_id"`
	TypeID 		primitive.ObjectID `json:"type" bson:"type"`
	AuthorID 	string `json:"author_id" bson:"author_id"`
	CreatedAt  primitive.DateTime `json:"created_at" bson:"created_at"`
}

type ProjectIncoming struct {
	Title    string `json:"title" bson:"title"`
	Content  string `json:"content" bson:"content"`
	Link	 string `json:"link" bson:"link"`
	Tags	 string `json:"tags" bson:"tags"`
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
	Comments int `json:"comments" bson:"comments,omitempty"`
	Awards []Award `json:"awards" bson:"awards,omitempty"`
	AwardsTotal int `json:"awards_total" bson:"awards_total,omitempty"`
	CreatedAt string `json:"created_at" bson:"created_at"`
}

type ProjectVote struct {
	ProjectID primitive.ObjectID `json:"project_id" bson:"project_id"`
	AuthorID  string             `json:"author_id" bson:"author_id"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
}