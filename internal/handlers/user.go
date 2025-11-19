package handlers

import (
	"full-domain/internal/services"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SignupHandler(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")
		password := c.PostForm("password")
		if userService.CreateUser(name, email, password) != nil {
			c.String(http.StatusInternalServerError, "Error")
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}

func LoginHandler(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")
		user, err := userService.Authenticate(email, password)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error")
			return
		}
		now := time.Now()
		user.LastLogin = &now

		if userService.Update(user) != nil {
			c.String(http.StatusInternalServerError, "Error")
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
