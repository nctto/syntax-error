package project

import (
	aw "go-api/internal/award"
	cm "go-api/internal/comment"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Project struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title    string             `json:"title" bson:"title"`
	Content  string             `json:"content" bson:"content"`
	AuthorID string 			`json:"author_id" bson:"author_id,omitempty"`
	Link	 string             `json:"link" bson:"link"`
	Tags	 []string            `json:"tags" bson:"tags,omitempty"`
	VotesTotal int32            	`json:"votes_total" bson:"votes_total,omitempty"`
	Voted	 bool            	`json:"voted" bson:"voted,omitempty"`
	CreatedAt primitive.DateTime `json:"created_at" bson:"created_at"`
	CommentsTotal int 			 `json:"comments_total" bson:"comments_total,omitempty"`
	Comments []cm.Comment 		`json:"comments" bson:"comments,omitempty"`
	Awards  []aw.Award 			`json:"awards" bson:"awards,omitempty"`
	AwardsTotal  int 			`json:"awards_total" bson:"awards_total,omitempty"`
}

type ProjectView struct {
	ID      string `json:"id" bson:"_id,omitempty"`
	Title    string `json:"title" bson:"title"`
	Content  string `json:"content" bson:"content"`
	AuthorID string `json:"author_id" bson:"author_id,omitempty"`
	Link	 string `json:"link" bson:"link"`
	Tags	 []string `json:"tags" bson:"tags,omitempty"`
	VotesTotal	 int32 `json:"votes_total" bson:"votes_total,omitempty"`
	Voted	 bool  `json:"voted" bson:"voted,omitempty"`
	CommentsTotal int `json:"comments_total" bson:"comments_total,omitempty"`
	Comments []cm.CommentView `json:"comments" bson:"comments,omitempty"`
	Awards []aw.Award `json:"awards" bson:"awards,omitempty"`
	AwardsTotal int `json:"awards_total" bson:"awards_total,omitempty"`
	CreatedAt string `json:"created_at" bson:"created_at"`
	TargetID string `json:"target_id" bson:"target_id"`
}

type ProjectIncoming struct {
	Title    string `json:"title" bson:"title"`
	Content  string `json:"content" bson:"content"`
	Link	 string `json:"link" bson:"link"`
	Tags	 string `json:"tags" bson:"tags"`
}