package main

import (
	"encoding/json"
	"go-api/cmd/server"
	"net/http"
	"testing"

	"go-api/internal/project"

	"github.com/alecthomas/assert/v2"
	"github.com/appleboy/gofight/v2"
	"github.com/gin-gonic/gin"
)




func init() {
	gin.SetMode(gin.TestMode)
}

func TestProjects(t *testing.T) {
	g := gofight.New()
	e := gin.Default()
	server.Run()

	var projectID string
	var basePath = "/projects"

	t.Run("GetProjects", func(t *testing.T) {
		g.GET(basePath).Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
	})
	
	t.Run("GetRandomProject", func(t *testing.T) {
		g.GET(basePath+"/random").Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
			var p project.Project
			json.Unmarshal(r.Body.Bytes(), &p)
			projectID = p.ID.Hex()
		})
	})

	t.Run("CreateProject", func(t *testing.T) {
		body := project.FakeProject()
		g.POST(basePath).SetBody(body).Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusCreated, r.Code)
			var p project.Project
			json.Unmarshal(r.Body.Bytes(), &p)
		})
	})

	t.Run("GetProjectByID", func(t *testing.T) {
		g.GET(basePath+ "/"+ projectID ).Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
	})

	t.Run("UpdateProject", func(t *testing.T) {
		body := `{
			"name": "John - Updated",
			"gender": "Does"
		}`
		g.PUT(basePath+"/"+ projectID).SetBody(body).Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
	})

	t.Run("DeleteProject", func(t *testing.T) {
		g.DELETE(basePath + "/"+ projectID ).Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
	})
}