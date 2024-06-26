package project

import (
	"context"
	"go-api/pkg/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var projectCollection = db.GetCollection("projects")

func DbGetAllProjects(page int, limit int, sortBy string, user interface{} ) ([]Project, int64, error) {

	var projects []Project
	pipeline := GetProjectsPaginatedPipeline(page, limit, sortBy)
	
	if user != nil {
		nickname := user.(map[string]interface{})["nickname"].(string)
		if nickname != "" {
			pipeline = AddProjectsVotedPipeline(pipeline, nickname)
		}
	}
	pipeline = AddProjectsPipelineSorter(pipeline, sortBy)

	cursor, err := projectCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return projects, 0, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var project Project
		cursor.Decode(&project)
		projects = append(projects, project)
	}

	totalRecords, err := projectCollection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return projects, 0, err
	}

	return projects, totalRecords, nil
}

func DbGetProjectsByUser(username string) ([]Project, error) {
	var projects []Project
	pipeline := GetProjectsByUserPipeline(username)
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

func DbGetProjectID(id primitive.ObjectID, user interface{}) (Project, error) {
	
	pipeline := GetProjectPipelineByID(id)
	if user != nil {
		nickname := user.(map[string]interface{})["nickname"].(string)
		pipeline = AddProjectsVotedPipeline(pipeline, nickname)
	}

	var project Project
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

func DbProjectExists(id primitive.ObjectID) (bool, error) {
	count, err := projectCollection.CountDocuments(context.Background(), bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func DbCreateProject(project Project) (primitive.ObjectID, error) {
	project.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
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