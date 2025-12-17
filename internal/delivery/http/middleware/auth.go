package middleware

import (
	"net/http"

	"full-domain/pkg/woodpecker"
	"full-domain/internal/domain"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired(userService domain.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := woodpecker.FromContext(c.Request.Context())

		session := sessions.DefaultMany(c, "user_session")
		email := session.Get("email")

		if email == nil {
			logger.Warn("unauthenticated request", "path", c.FullPath())
			session.Clear()
			session.Save()
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		user, err := userService.FindByEmail(email.(string))
		if err != nil || user == nil {
			logger.Warn("auth user not found", "email", email)
			session.Clear()
			session.Save()
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		c.Next()
	}
}
