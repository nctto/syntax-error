package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-api/cmd/server/middleware"

	"github.com/gin-contrib/sessions"

	usr "go-api/internal/user"
	wlt "go-api/internal/wallet"
)

func InitializeCreatePage(router *gin.Engine) {
	router.GET("/create", middleware.IsAuthenticated, func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("profile")
		if user == nil {
			c.Redirect(http.StatusFound, "/auth/login")
		}
		nickname := usr.UserNickName(user)
		balance, err := wlt.DbGetUserWalletBalanceByNickName(nickname)
		if err != nil {
			c.HTML(200, "create-project-page.html", gin.H{
				"title": "syntax error", 
				"session_user": user,
				"balance": "XXXX",
			})
		}
		c.HTML(200, "create-project-page.html", gin.H{
			"title": "syntax error", 
			"session_user": user,
			"balance": balance,
		})
})}