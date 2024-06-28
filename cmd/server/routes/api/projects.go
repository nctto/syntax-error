package api

import (
	"fmt"
	"go-api/cmd/server/middleware"
	"go-api/internal/project"
	"html/template"

	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleSocket(c *gin.Context) {
	
	session := sessions.Default(c)
	user := session.Get("profile")
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}


	roomId := c.Param("id")
	fmt.Printf("Room: %s\n", roomId)
	room, err := primitive.ObjectIDFromHex(roomId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid room ID"})
		return
	}

	fmt.Printf("Room: %s\n", room)

	
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		// panic(err)
		fmt.Printf("%s, error while Upgrading websocket connection\n", err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	for {
		// Read message from client
		messageType, p, err := conn.ReadMessage()
		fmt.Printf("New Message: %s\n %d\n", p, messageType)
		if err != nil {
			// panic(err)
			fmt.Printf("%s, error while reading message\n", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
			break
		}

		// parse html string
		html := template.HTMLEscapeString("<h1> Hello World "+string(p)+"</h1>")
		

		fmt.Printf("HTML: %s\n", html)

		toBytes := []byte(html)

		// Echo message back to client
		err = conn.WriteMessage(messageType, toBytes)
		if err != nil {
			// panic(err)
			fmt.Printf("%s, error while writing message\n", err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)
			break
		}
	}
}

func InitializeProjects(router *gin.Engine) {
	router.GET("/api/projects", middleware.IsAuthenticated, project.GetProjects)
	router.POST("/api/projects", middleware.IsAuthenticated, project.CreateProject)
	router.GET("/api/projects/:id", middleware.IsAuthenticated, project.GetProjectByID)
	router.PUT("/api/projects/:id",  middleware.IsAuthenticated, project.UpdateProject)
	router.DELETE("/api/projects/:id", middleware.IsAuthenticated, project.DeleteProject)
	router.GET("/api/projects/fake",middleware.IsAuthenticated, project.FakeProjects)
	// router.GET("/api/projects/random", middleware.IsAuthenticated, project.GetRandomProject)

	router.GET("/ws/:id", middleware.IsAuthenticated, HandleSocket)



	router.POST("/api/upload", middleware.IsAuthenticated, project.UploadFile)


}