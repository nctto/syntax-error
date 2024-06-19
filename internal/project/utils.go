package project

import (
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
	if project.AuthorID == "" {
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
		TargetID: ObjectIdToString(project.ID),
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
		Comments: cm.CommentsPaginatedView(project.ID, project.Comments, int64(project.CommentsTotal), 10, 1, "best"),
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

func ProjectPaginatedView(projects []Project, totalRecords int64, page int, limit int, sortBy string) ProjectPaginated {
	ProjectPaginated := ProjectPaginated{}
	ProjectPaginated.Data = ProjectsToProjectView(projects)

	pagination := Pagination{}
	pagination.Page = page
	pagination.Limit = limit
	pagination.SortBy = sortBy
	pagination.TotalPages = totalRecords / int64(limit)
	pagination.TotalRecords = totalRecords
	pagination.CurrentPage = int64(page)
	if page < int(pagination.TotalPages) {
		pagination.HasNext = true
	} else {
		pagination.HasNext = false
	}

	if page > 0 {
		pagination.HasPrev = true
	} else {
		pagination.HasPrev = false
	}
	pagination.NextLink = "/api/projects?page=" + strconv.Itoa(page+1) + "&limit=" + strconv.Itoa(limit)
	pagination.PrevLink = "/api/projects?page=" + strconv.Itoa(page-1) + "&limit=" + strconv.Itoa(limit)
	ProjectPaginated.Pagination = pagination

	return ProjectPaginated
}

func AddProjectsPipelineSorter(pipeline []bson.M, sortBy string) []bson.M {
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
	return p, l, sortBy
}

func GetPipeline(pipeline []bson.M) []bson.M {
	dflt := []bson.M{
		{"$lookup": bson.M{
			"from":         "votes",
			"localField":   "_id",
			"foreignField": "target_id",
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

	pipeline = append(pipeline, dflt...)
	return pipeline
}

func GetProjectPipelineByID(projectID primitive.ObjectID) []bson.M {
	pipeline := []bson.M{
		{"$match": bson.M{"_id": projectID}},
	}
	pipeline = GetPipeline(pipeline)
	return pipeline
}
func GetProjectsByUserPipeline(username string) []bson.M {
	pipeline := []bson.M{
		{"$match": bson.M{"author_id": username}},
	}
	pipeline = GetPipeline(pipeline)
	return pipeline
}

func GetProjectsPaginatedPipeline(page int, limit int, sortBy string) []bson.M {
	skip := int64(page*limit - limit)
	
	// lookup most voted projects
	if sortBy == "best" {
		pipeline := []bson.M{
			{
				"$addFields": bson.M{
					"totalRecords": "$totalRecords.count",
				},
			},

			{
				"$lookup": bson.M{
					"from":         "votes",
					"localField":   "_id",
					"foreignField": "target_id",
					"as":           "votes",
				},
			},
			{
				"$addFields": bson.M{
					"votes_total": bson.M{"$size": "$votes"},
				},
			},
			{
				"$sort": bson.M{"votes_total": -1},
			},
			{
				"$skip": skip,
			},
			{
				"$limit": limit,
			},
		};
		pipeline = GetPipeline(pipeline)
		return pipeline
	}
	pipeline := []bson.M{
		{"$skip": skip},
		{"$limit": limit},
	}
	pipeline = GetPipeline(pipeline)
	return pipeline
}

func AddProjectsVotedPipeline(pipeline []bson.M, authorID string) []bson.M {
	pipeline = append(pipeline, bson.M{"$lookup": bson.M{
		"from":         "votes",
		"let":          bson.M{"target_id": "$_id"},
		"pipeline":     []bson.M{
			{"$match": bson.M{
				"$expr": bson.M{
					"$and": []bson.M{
						{"$eq": []string{"$target_id", "$$target_id"}},
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