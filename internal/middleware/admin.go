package middleware

import (
	"full-domain/internal/lumberjack"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.DefaultMany(c, "admin_session")

		role := session.Get("role")
		if role == nil || role.(string) != "admin" {
			lumberjack.Logger.Warn("unauthorized admin access", "path", c.FullPath())
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		c.Next()
	}
}
