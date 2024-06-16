package comment

import (
	"context"
	"go-api/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var commentCollection = db.GetCollection("comments")

func DbGetAllComments(page int, limit int, sortBy string, user interface{} , targetID primitive.ObjectID) ([]Comment, error) {

	var comments []Comment
	pipeline := GetCommentsPipeline(page, limit, sortBy, user, targetID)
	if user != nil {
		nickname := user.(map[string]interface{})["nickname"].(string)
		
		if nickname != "" {
			pipeline = AddCommentsVotedPipeline(pipeline, nickname)
		}
	}
			
	pipeline = AddCommentsPipelineSorter(pipeline, sortBy)	
	cursor, err := commentCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return comments, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var comment Comment
		cursor.Decode(&comment)
		comment.Replies, err = DbGetAllComments(1, 10, "new", user, comment.ID)
		if err != nil {
			comment.Replies = []Comment{}
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func DbCommentExists(id primitive.ObjectID) (bool, error) {
	count, err := commentCollection.CountDocuments(context.Background(), bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func DbGetCommentID(id primitive.ObjectID) (Comment, error) {
	var comment Comment
	err := commentCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&comment)
	return comment, err
}

func DbCreateComment(comment Comment) (primitive.ObjectID, error) {
	result, err := commentCollection.InsertOne(context.Background(), comment)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func DbUpdateComment(id primitive.ObjectID, comment Comment) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": comment}
	_, err := commentCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func DbDeleteComment(id primitive.ObjectID) error {
	_, err := commentCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
