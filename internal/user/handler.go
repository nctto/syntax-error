package user

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUsers(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")

	users, err := DbGetAllUsers(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid  ID"})
		return
	}

	user, err := DbGetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := DbCreateUser(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	var updatedUser User
	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = DbUpdateUser(id, updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	updatedUser.ID = id
	c.JSON(http.StatusOK, updatedUser)
}

func DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID"})
		return
	}

	err = DbDeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
