package vote

import (
	"context"
	"fmt"
	"go-api/pkg/db"

	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var voteCollection = db.GetCollection("votes")

func DbGetAllVotes(page string, limit string) ([]Vote, error) {
	var votes []Vote

	l, _ := strconv.ParseInt(limit, 10, 64)
	p, _ := strconv.ParseInt(page, 10, 64)

	skip := int64(p*l - l)

	fOpt := options.FindOptions{
		Skip:  &skip,
		Limit: &l,
	}

	cursor, err := voteCollection.Find(context.Background(), bson.M{}, &fOpt)

	if err != nil {
		return votes, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var vote Vote
		cursor.Decode(&vote)
		votes = append(votes, vote)
	}

	return votes, nil
}

func DbGetVoteID(id primitive.ObjectID) (Vote, error) {
	var vote Vote
	err := voteCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&vote)
	return vote, err
}

func DbVoteExists(projectID primitive.ObjectID, user string) (bool, error) {
	count, err := voteCollection.CountDocuments(context.Background(), bson.M{"project_id": projectID, "author_id": user})
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return count > 0, nil
}


func DbCreateVote(vote Vote) (primitive.ObjectID, error) {
	result, err := voteCollection.InsertOne(context.Background(), vote)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func DbUpdateVote(id primitive.ObjectID, vote Vote) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": vote}
	_, err := voteCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func DbDeleteVote(id primitive.ObjectID) error {
	_, err := voteCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func DbDeleteVoteByAuthor(projectID primitive.ObjectID, authorID string) error {
	_, err := voteCollection.DeleteOne(context.Background(), bson.M{"project_id": projectID, "author_id": authorID})
	return err
}


func DbGetProjectVotes(projectID primitive.ObjectID) (int32, error) {
	var votes []Vote
	cursor, err := voteCollection.Find(context.Background(), bson.M{"project_id": projectID})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var vote Vote
		cursor.Decode(&vote)
		votes = append(votes, vote)
	}
	return int32(len(votes)), nil	
}