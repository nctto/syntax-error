package comment

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RequiredFields(comment Comment) bool {
	if comment.TargetID == primitive.NilObjectID {
		return false
	}
	if comment.AuthorID == "" {
		return false
	
	}
	if comment.Content == "" {
		return false
	}
	return true
}

func ObjectIdToString(id primitive.ObjectID) string {
	return id.Hex()
}

func DateToString(date primitive.DateTime) string {
	return date.Time().String()
}

func AddCommentsPipelineSorter(pipeline []bson.M, sortBy string) []bson.M {
	fmt.Println("Sort by", sortBy)
	if sortBy == "new"{
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"created_at": -1}})
	} else if sortBy == "old"{
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"created_at": 1}})
	} else if sortBy == "uncommentd"{
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"comments": 1}})
	} else {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"comments": -1}})
	}
	return pipeline
}

func CommentsDefaultQueryParams(c *gin.Context) (int, int, string) {
	page := c.Query("page")
	limit := c.Query("limit")
	sortBy := c.Query("sort_by")
	var p, l int = 0, 0
	if page == "" {
		p = 1
	} else {
		p,_ = strconv.Atoi(page)
	}
	if limit == "" {
		l = 10
	} else {
		l,_ = strconv.Atoi(limit)
	}
	if sortBy == "" {
		sortBy = "best"
	}

	if l > 100 {
		l = 100
	}
	fmt.Println("Page", p, "Limit", l, "Sort by", sortBy)
	return p, l, sortBy
}

func GetCommentsPipeline(page int, limit int, sortBy string) []bson.M {
	skip := int64(page*limit - limit)
	pipeline := []bson.M{
		{"$skip": skip},
		{"$limit": limit},
	}
	return pipeline
}