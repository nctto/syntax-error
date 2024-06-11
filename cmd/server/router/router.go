package router

import (
	"encoding/gob"
	"net/http"
	"time"

	"go-api/cmd/server/authenticator"

	api "go-api/cmd/server/routes/api"
	pages "go-api/cmd/server/routes/pages"
	"go-api/cmd/server/routes/ui"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"go-api/cmd/server/routes"
)

const userkey = "user"

var secret = []byte("secret")

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.Next()
}

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, rate ratelimit.Info) {
	c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
}

func New(auth *authenticator.Authenticator) *gin.Engine {

	router := gin.Default()
	
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{Rate: time.Second, Limit: 5})
	mw := ratelimit.RateLimiter(store, &ratelimit.Options{ErrorHandler: errorHandler, KeyFunc: keyFunc})
	
	gob.Register(map[string]interface{}{})
	
	authStore := cookie.NewStore([]byte("secret"))
	
	
	router.Use(mw)
	router.Use(sessions.Sessions("auth-session", authStore))
	

	router.Static("/assets", "/Users/rfcku/sites/rfcku/go-api/cmd/server/assets")
	router.LoadHTMLGlob("/Users/rfcku/sites/rfcku/go-api/cmd/server/html/templates/**/*")
	
	routes.InitializeAuth(router, auth)
	api.InitializeProjects(router)
	api.InitializeVotes(router)
	ui.InitializeCommentsUI(router)
	ui.InitializeVotesUI(router)
	ui.InitializeProjectsUI(router)
	ui.InitializeLikesUI(router)
	pages.InitializeCreatePage(router)
	pages.InitializeHomePage(router)
	pages.InitializeSingleProjectPage(router)
	return router
}