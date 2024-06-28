package project

import (
	"encoding/json"
	"fmt"
	usr "go-api/internal/user"
	wlt "go-api/internal/wallet"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
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

	session := sessions.Default(c)
	user := session.Get("profile")
	
	page, limit, sortBy := ProjectsDefaultQueryParams(c)
	projects, total, err := DbGetAllProjects(page, limit, sortBy, user)
	
	paginatedProjects := ProjectPaginatedView(projects, total, page, limit, sortBy)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, paginatedProjects)
}

func GetProjectByID(c *gin.Context) {
	
	session := sessions.Default(c)
	user := session.Get("profile")
	
	
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
		return
		}
	fmt.Println("Getting Project by ID", id, user)

	project, err := DbGetProjectID(id, user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Project not found"})
		return
	}
	projectView := ProjectToProjectView(project)
	c.JSON(http.StatusOK, projectView)
}

func CreateProject(c *gin.Context) {

	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var project = ProjectIncoming{}
	if err := c.BindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var newProject = Project{
		Title: project.Title,
		Content: project.Content,
		Link: project.Link,
		Tags: strings.Split(project.Tags, ","),
	}

	newProject.AuthorID = user.(map[string]interface{})["nickname"].(string)
	newProject.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
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

func GetHTMLCreateForm(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.Redirect(http.StatusFound, "/auth/login")
	}
	nickname := usr.UserNickName(user)
	balance, err := wlt.DbGetUserWalletBalanceByNickName(nickname)
	if err != nil {
		c.HTML(200, "create-project-page.html", gin.H{
			"title": "Submit Projec", 
			"session_user": user,
			"balance": "XXXX",
		})
	}
	c.HTML(200, "create-project-page.html", gin.H{
		"title": "Submit Project", 
		"session_user": user,
		"balance": balance,
	})
}

func GetHTMLAllProjects(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("profile")

	page, limit, sortBy := ProjectsDefaultQueryParams(c)
	projects, total,  err := DbGetAllProjects(page, limit, sortBy, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	paginatedProjects := ProjectPaginatedView(projects, total, page, limit, sortBy)

	c.HTML(200, "list-projects.html", gin.H{
		"session_user": user,
		"projects": paginatedProjects,
	})
}

func GetHTMLSubmitProjectForm(c *gin.Context) {
		
	session := sessions.Default(c)
	user := session.Get("profile")

	if user == nil {
		c.Redirect(http.StatusFound, "/auth/login")
	}

	var errors = []gin.H{}
	var newProject Project
	if err := c.BindJSON(&newProject); err != nil {
		errors = append(errors, gin.H{"message": err.Error()})
	}
	newProject.AuthorID = user.(map[string]interface{})["nickname"].(string)
	if !RequiredFields(newProject) {
		errors = append(errors, gin.H{"message": "Missing required fields"})
	}

	if len(errors) > 0 {
		c.HTML(200, "response-create-project.html", gin.H{
			"submitted": false,
			"message": "Error creating project",
			"errors": errors,
		})
		return
	}

	id, err := DbCreateProject(newProject)
	if err != nil {
		errors = append(errors, gin.H{"message": err.Error()})
	}
	newProject.ID = id
	c.HTML(200, "response-create-project.html", gin.H{
		"submitted": true,
		"message": "Project created successfully",
		"errors": errors,
	})
}

func UploadFile(c *gin.Context) {
		
		session := sessions.Default(c)
		user := session.Get("profile")
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		// Multipart form
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, "get form err: %s", err.Error())
			return
		}
		files := form.File["files"]

		for _, file := range files {
			filename := filepath.Base(file.Filename)
			if err := c.SaveUploadedFile(file, filename); err != nil {
				c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
				return
			}
		}

		c.String(http.StatusOK, "Uploaded successfully %d files with fields name=%s and email=%s.", len(files))
}