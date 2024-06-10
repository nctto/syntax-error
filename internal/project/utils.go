package project

import (
	"fmt"
	cm "go-api/internal/comment"
	utils "go-api/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RequiredFields(project Project) bool {
	if project.Title == "" {
		return false
	}
	if project.Content == "" {
		return false
	}
	if project.Link == "" {
		return false
	}
	return true
}

func fakeProject() string {
	return `{
		"title": "`+faker.Word()+`",
		"content": "`+faker.Sentence()+`",
		"link": "`+faker.URL()+`",
		"author_id": "`+faker.Username()+`",
		"tags": ["`+faker.Word()+`", "`+faker.Word()+`"]
	}`
}

func ObjectIdToString(id primitive.ObjectID) string {
	return id.Hex()
}

func ProjectToProjectView(project Project) ProjectView {
	return ProjectView{
		ID: ObjectIdToString(project.ID),
		Title: project.Title,
		Content: project.Content,
		AuthorID: project.AuthorID,
		Link: project.Link,
		Tags: project.Tags,
		VotesTotal: project.VotesTotal,
		Voted: project.Voted,
		CommentsTotal: project.CommentsTotal,
		Awards: project.Awards,
		AwardsTotal: project.AwardsTotal,
		Comments: cm.CommentsToCommentView(project.Comments),
		CreatedAt: utils.DateToString(project.CreatedAt),

	}
}

func ProjectsToProjectView(projects []Project) []ProjectView {
	var projectView []ProjectView
	for _, project := range projects {
		projectView = append(projectView, ProjectToProjectView(project))
	}
	return projectView
}

func AddProjectsPipelineSorter(pipeline []bson.M, sortBy string) []bson.M {
	fmt.Println("Sort by", sortBy)
	if sortBy == "new"{
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"created_at": -1}})
	} else if sortBy == "old"{
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"created_at": 1}})
	} else if sortBy == "unvoted"{
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"votes": 1}})
	} else {
		pipeline = append(pipeline, bson.M{"$sort": bson.M{"votes": -1}})
	}
	return pipeline
}

func ProjectsDefaultQueryParams(c *gin.Context) (int, int, string) {
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

func GetProjectPipeline(projectID primitive.ObjectID) []bson.M {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": projectID}},
		{"$lookup": bson.M{
			"from":         "votes",
			"localField":   "_id",
			"foreignField": "project_id",
			"as":           "votes",
		}},
		{"$lookup": bson.M{
			"from":         "comments",
			"localField":   "_id",
			"foreignField": "target_id",
			"as":           "comments",
		}},
		{"$lookup": bson.M{
			"from":         "awards",
			"localField":   "_id",
			"foreignField": "target_id",
			"as":           "awards",
		}},
		{"$addFields": bson.M{
			"votes_total": bson.M{"$size": "$votes"},
		}},
		{"$addFields": bson.M{
			"comments_total": bson.M{"$size": "$comments"},
		}},
		{"$addFields": bson.M{
			"awards_total": bson.M{"$size": "$awards"},
		}},
		{"$addFields": bson.M{
			// total awards and first 5 awards
			"awards": bson.M{"$slice": []interface{}{"$awards", 3}},
		}},
	}
	return pipeline

}


func GetProjectsPipeline(page int, limit int, sortBy string) []bson.M {
	skip := int64(page*limit - limit)
	pipeline := []bson.M{
		{"$skip": skip},
		{"$limit": limit},
		{"$lookup": bson.M{
			"from":         "votes",
			"localField":   "_id",
			"foreignField": "project_id",
			"as":           "votes",
		}},
		{"$lookup": bson.M{
			"from":         "comments",
			"localField":   "_id",
			"foreignField": "target_id",
			"as":           "comments",
		}},
		{"$lookup": bson.M{
			"from":         "awards",
			"localField":   "_id",
			"foreignField": "target_id",
			"as":           "awards",
		}},
		{"$addFields": bson.M{
			"votes_total": bson.M{"$size": "$votes"},
		}},
		{"$addFields": bson.M{
			"comments_total": bson.M{"$size": "$comments"},
		}},
		{"$addFields": bson.M{
			"awards_total": bson.M{"$size": "$awards"},
		}},
		{"$addFields": bson.M{
			// total awards and first 5 awards
			"awards": bson.M{"$slice": []interface{}{"$awards", 3}},
		}},
	}
	return pipeline
}

func AddProjectsVotedPipeline(pipeline []bson.M, authorID string) []bson.M {
	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "votes",
		"let":          bson.M{"project_id": "$_id"},
		"pipeline":     []bson.M{
			{"$match": bson.M{
				"$expr": bson.M{
					"$and": []bson.M{
						{"$eq": []string{"$project_id", "$$project_id"}},
						{"$eq": []string{"$author_id", authorID}},
					},
				},
			}},
		},
		"as": "voted",
	}})
	pipeline = append(pipeline, bson.M{
		"$addFields": bson.M{
			"voted": bson.M{"$cond": bson.M{
				"if":  bson.M{
					"$gt": []interface{}{
						bson.M{"$size": "$voted"},
						0,
					},
				},
				"then": true,
				"else": false,
			}},
		},
	})
	return pipeline
}