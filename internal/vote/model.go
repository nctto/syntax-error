package vote

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)



type Vote struct {
	ID      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProjectID primitive.ObjectID `json:"project_id" bson:"project_id"`
	UserID primitive.ObjectID `json:"user_id" bson:"user_id"`
	Vote	 int            	`json:"vote" bson:"vote"`
	CreatedAt string `json:"created_at" bson:"created_at"`
}