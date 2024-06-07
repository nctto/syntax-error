package vote

import (
	"context"
	"go-api/pkg/db"

	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var projectCollection = db.GetCollection("projects")

func DbGetAllVotes(page string, limit string) ([]Vote, error) {
	var projects []Vote

	l, _ := strconv.ParseInt(limit, 10, 64)
	p, _ := strconv.ParseInt(page, 10, 64)

	skip := int64(p*l - l)

	fOpt := options.FindOptions{
		Skip:  &skip,
		Limit: &l,
	}

	cursor, err := projectCollection.Find(context.Background(), bson.M{}, &fOpt)

	if err != nil {
		return projects, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var project Vote
		cursor.Decode(&project)
		projects = append(projects, project)
	}

	return projects, nil
}

func DbGetVoteID(id primitive.ObjectID) (Vote, error) {
	var project Vote
	err := projectCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&project)
	return project, err
}

func DbCreateVote(project Vote) (primitive.ObjectID, error) {
	result, err := projectCollection.InsertOne(context.Background(), project)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func DbUpdateVote(id primitive.ObjectID, project Vote) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": project}
	_, err := projectCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func DbDeleteVote(id primitive.ObjectID) error {
	_, err := projectCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}