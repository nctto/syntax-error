package like

import (
	"context"
	"fmt"
	"go-api/pkg/db"

	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var likeCollection = db.GetCollection("likes")

func DbGetAllLikes(page string, limit string) ([]Like, error) {
	var likes []Like

	l, _ := strconv.ParseInt(limit, 10, 64)
	p, _ := strconv.ParseInt(page, 10, 64)

	skip := int64(p*l - l)

	fOpt := options.FindOptions{
		Skip:  &skip,
		Limit: &l,
	}

	cursor, err := likeCollection.Find(context.Background(), bson.M{}, &fOpt)

	if err != nil {
		return likes, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var like Like
		cursor.Decode(&like)
		likes = append(likes, like)
	}

	return likes, nil
}

func DbGetLikeID(id primitive.ObjectID) (Like, error) {
	var like Like
	err := likeCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&like)
	return like, err
}

func DbLikeExists(commentID primitive.ObjectID, user string) (bool, error) {
	count, err := likeCollection.CountDocuments(context.Background(), bson.M{"comment_id": commentID, "author_id": user})
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return count > 0, nil
}


func DbCreateLike(like Like) (primitive.ObjectID, error) {
	result, err := likeCollection.InsertOne(context.Background(), like)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func DbUpdateLike(id primitive.ObjectID, like Like) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": like}
	_, err := likeCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func DbDeleteLike(id primitive.ObjectID) error {
	_, err := likeCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func DbDeleteLikeByAuthor(commentID primitive.ObjectID, authorID string) error {
	_, err := likeCollection.DeleteOne(context.Background(), bson.M{"comment_id": commentID, "author_id": authorID})
	return err
}


func DbGetCommentLikes(commentID primitive.ObjectID) (int32, error) {
	var likes []Like
	cursor, err := likeCollection.Find(context.Background(), bson.M{"comment_id": commentID})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var like Like
		cursor.Decode(&like)
		likes = append(likes, like)
	}
	return int32(len(likes)), nil	
}