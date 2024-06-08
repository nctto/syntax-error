package project

import (
	"context"
	"go-api/pkg/db"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var projectCollection = db.GetCollection("projects")

func DbGetAllProjects(page string, limit string) ([]Project, error) {
	var projects []Project

	l, _ := strconv.ParseInt(limit, 10, 64)
	p, _ := strconv.ParseInt(page, 10, 64)

	skip := int64(p*l - l)

	pipeline := []bson.M{
		{"$sort": bson.M{"created_at": -1}},
		{"$skip": skip},
		{"$limit": l},
		{"$lookup": bson.M{
			"from":         "votes",
			"localField":   "_id",
			"foreignField": "project_id",
			"as":           "votes",
		}},
		{"$addFields": bson.M{
			"votes": bson.M{"$size": "$votes"},
		}},
		{"$project": bson.M{
			"id":    "$_id",
			"title": "$title",
			"content": "$content",
			"author_id": "$author_id",
			"link": "$link",
			"tags": "$tags",
			"created_at": "$created_at",
			"votes": "$votes",
		}},
		{"$sort": bson.M{"votes": -1}},
	}

	cursor, err := projectCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return projects, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var project Project
		cursor.Decode(&project)
		projects = append(projects, project)
	}

	return projects, nil
}

func DbGetProjectID(id primitive.ObjectID) (Project, error) {
	var project Project
	err := projectCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&project)
	return project, err
}

func DbCreateProject(project Project) (primitive.ObjectID, error) {
	result, err := projectCollection.InsertOne(context.Background(), project)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func DbUpdateProject(id primitive.ObjectID, project Project) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": project}
	_, err := projectCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func DbDeleteProject(id primitive.ObjectID) error {
	_, err := projectCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func DbGetRandomProject() (Project, error) {
	var project Project
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
