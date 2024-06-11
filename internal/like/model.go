package like

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Like struct {
	ID      	primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TargetID 	primitive.ObjectID `json:"target_id" bson:"target_id"`
	AuthorID 	string `json:"author_id" bson:"author_id"`
	Liked	 	bool            	`json:"liked" bson:"liked,omitempty"`
	CreatedAt  primitive.DateTime `json:"created_at" bson:"created_at"`
}