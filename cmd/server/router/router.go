package router

import (
	"encoding/gob"

	"go-api/cmd/server/authenticator"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"go-api/cmd/server/routes"
)


func New(auth *authenticator.Authenticator) *gin.Engine {

	router := gin.Default()
	
	// store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{Rate: time.Second, Limit: 5})
	// mw := ratelimit.RateLimiter(store, &ratelimit.Options{ErrorHandler: errorHandler, KeyFunc: keyFunc})
	
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	router.Static("/assets", "/Users/rfcku/sites/rfcku/go-api/cmd/server/assets")
	router.LoadHTMLGlob("/Users/rfcku/sites/rfcku/go-api/cmd/server/templates/*")
	
	routes.InitializeAuth(router, auth)
	routes.InitializeHome(router, auth)
	routes.InitializeProjects(router, auth)
	return router
}