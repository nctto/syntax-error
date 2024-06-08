package vote

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)



type Vote struct {
	ID      	primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProjectID 	primitive.ObjectID `json:"project_id" bson:"project_id"`
	AuthorID 	string `json:"author_id" bson:"author_id"`
	Vote	 int            	`json:"vote" bson:"vote"`
	CreatedAt  primitive.DateTime `json:"created_at" bson:"created_at"`
}