package project

import (
	"github.com/go-faker/faker/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RequiredFields(newProject Project) bool {
	if newProject.Title == "" {
		return false
	}
	if newProject.Content == "" {
		return false
	}
	if newProject.Link == "" {
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

func DateToString(date primitive.DateTime) string {
	return date.Time().String()
}

func ProjectToProjectView(project Project) ProjectView {
	return ProjectView{
		ID: ObjectIdToString(project.ID),
		Title: project.Title,
		Content: project.Content,
		AuthorID: project.AuthorID,
		Link: project.Link,
		Tags: project.Tags,
		Votes: project.Votes,
		CreatedAt: DateToString(project.CreatedAt),
	}

}

func ProjectsToProjectView(projects []Project) []ProjectView {
	var projectView []ProjectView
	for _, project := range projects {
		projectView = append(projectView, ProjectToProjectView(project))
	}
	return projectView
}