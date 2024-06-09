package comment

import (
	"context"
	"go-api/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var projectCollection = db.GetCollection("projects")

func DbGetAllComments(page int, limit int, sortBy string, user interface{} ) ([]Comment, error) {

	var projects []Comment
	pipeline := GetCommentsPipeline(page, limit, sortBy)
	pipeline = AddCommentsPipelineSorter(pipeline, sortBy)

	cursor, err := projectCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return projects, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var project Comment
		cursor.Decode(&project)
		projects = append(projects, project)
	}

	return projects, nil
}

func DbGetCommentID(id primitive.ObjectID) (Comment, error) {
	var project Comment
	err := projectCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&project)
	return project, err
}

func DbCreateComment(project Comment) (primitive.ObjectID, error) {
	result, err := projectCollection.InsertOne(context.Background(), project)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func DbUpdateComment(id primitive.ObjectID, project Comment) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": project}
	_, err := projectCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func DbDeleteComment(id primitive.ObjectID) error {
	_, err := projectCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func DbGetRandomComment() (Comment, error) {
	var project Comment
	pipeline := []bson.M{
		{"$sample": bson.M{"size": 1}},
	}
	cursor, err := projectCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return project, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		cursor.Decode(&project)
	}
	return project, nil
}
