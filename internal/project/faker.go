package project

import (
	"github.com/go-faker/faker/v4"
)

func fakeProject() string {
	return `{
		"title": "`+faker.Word()+`",
		"content": "`+faker.Sentence()+`",
		"link": "`+faker.URL()+`",
		"tags": ["`+faker.Word()+`", "`+faker.Word()+`"]
	}`
}