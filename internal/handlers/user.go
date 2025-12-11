package handlers

import (
	"full-domain/internal/lumberjack"
	"full-domain/internal/services"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SignupHandler(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")
		password := c.PostForm("password")

		if strings.TrimSpace(password) == "" {
			c.String(http.StatusBadRequest, "Password cannot be empty.")
			return
		}

		lowercase := regexp.MustCompile(`[a-z]`)
		uppercase := regexp.MustCompile(`[A-Z]`)
		digit := regexp.MustCompile(`\d`)
		special := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>_\-]`)

		if len(password) < 8 ||
			!lowercase.MatchString(password) ||
			!uppercase.MatchString(password) ||
			!digit.MatchString(password) ||
			!special.MatchString(password) {
			c.String(http.StatusBadRequest,
				"Password must be at least 8 characters long and include uppercase, lowercase, digit, and special character.")
			return
		}

		if userService.CreateUser(name, email, password) != nil {
			c.String(http.StatusInternalServerError, "Error")
			c.Redirect(http.StatusFound, "/signup")
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}

func LoginHandler(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := lumberjack.FromContext(ctx)

		email := c.PostForm("email")
		logger.Info("user login attempt", "email", email)
		password := c.PostForm("password")
		user, err := userService.Authenticate(c.Request.Context(), email, password)
		if err != nil {
			// c.String(http.StatusInternalServerError, "Error")
			c.Redirect(http.StatusFound, "/")
			return
		}
		now := time.Now()
		user.LastLogin = &now

		if userService.Update(user) != nil {
			// c.String(http.StatusInternalServerError, "Error")
			c.Redirect(http.StatusFound, "/")
			return
		}

		session := sessions.DefaultMany(c, "user_session")
		session.Set("email", user.Email)
		session.Set("name", user.Name)
		session.Set("role", user.Role)
		session.Save()

		session.Save()

		c.Redirect(http.StatusFound, "/home")
	}
}

func LogoutHandler(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.DefaultMany(c, "user_session")
		session.Clear()
		session.Save()

		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		c.Redirect(http.StatusFound, "/")
	}
}

func UpdatePasswordHandler(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.DefaultMany(c, "user_session")

		email := session.Get("email")
		if email == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		oldPassword := c.PostForm("old-password")
		newPassword := c.PostForm("new-password")

		err := userService.UpdatePassword(email.(string), oldPassword, newPassword)
		if err != nil {
			return
		}

		c.Redirect(http.StatusFound, "/home")
	}
}
