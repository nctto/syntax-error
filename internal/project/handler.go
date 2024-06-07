package project

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FakeProject() string {
	project := fakeProject()
	return project
}

func FakeProjects(c *gin.Context) {

	num := c.Query("num")
	numInt, err := strconv.Atoi(num)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid number"})
		return
	}
	projects := []Project{}
	for i := 0; i < numInt; i++ {
		project := fakeProject()
		var p Project 
		_, err := json.Marshal(project)
		if err != nil {
			println(err.Error())
			continue
		}
		err = json.Unmarshal([]byte(project), &p)
		if err != nil {
			println(err.Error())
			continue
		}
		DbCreateProject(p)
	}
	c.JSON(http.StatusCreated, projects)
}

func GetProjects(c *gin.Context) {

	page := c.Query("page")
	limit := c.Query("limit")

	projects, err := DbGetAllProjects(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projects)
}

func GetProjectByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
		return
	}

	project, err := DbGetProjectID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Project not found"})
		return
	}
	c.JSON(http.StatusOK, project)
}

func CreateProject(c *gin.Context) {

	var newProject Project
	if err := c.BindJSON(&newProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if !RequiredFields(newProject) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required fields"})
		return
	}

	id, err := DbCreateProject(newProject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	newProject.ID = id
	c.JSON(http.StatusCreated, newProject)
}

func UpdateProject(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	var updatedProject Project
	if err := c.BindJSON(&updatedProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = DbUpdateProject(id, updatedProject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	updatedProject.ID = id
	c.JSON(http.StatusOK, updatedProject)
}

func DeleteProject(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	err = DbDeleteProject(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project deleted"})
}

func GetRandomProject(c *gin.Context) {
	project, err := DbGetRandomProject()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, project)
}
