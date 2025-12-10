package middleware

import (
	"net/http"

	"full-domain/internal/lumberjack"
	"full-domain/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {

		session := sessions.DefaultMany(c, "user_session")

		email := session.Get("email")
		if email == nil {
			lumberjack.Logger.Warn("unauthenticated request", "path", c.FullPath())
			session.Clear()
			session.Save()
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		user, err := userService.FindByEmail(email.(string))
		if err != nil || user == nil {
			session.Clear()
			session.Save()
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}
