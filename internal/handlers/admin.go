package handlers

import (
	"full-domain/internal/models"
	"full-domain/internal/services"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AdminLoginHandler(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")

		user, err := userService.Authenticate(email, password)
		if err != nil || user.Role != "admin" {
			c.String(http.StatusUnauthorized, "Invalid credentials")
			return
		}
		session := sessions.DefaultMany(c, "admin_session")
		session.Set("email", user.Email)
		session.Set("name", user.Name)
		session.Set("role", user.Role)
		session.Save()

		c.Redirect(http.StatusFound, "/admin/dashboard")
	}
}

func AdminDashboardHandler(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.Query("q")
		editID := c.Query("edit")
		create := c.Query("create")

		users, _ := userService.SearchUsers(search)

		var editUser *models.User
		if editID != "" {
			editUser, _ = userService.FindByIDString(editID)
		}

		c.HTML(http.StatusOK, "admin-dashboard.html", gin.H{
			"users":      users,
			"search":     search,
			"editUser":   editUser,
			"createUser": create == "true",
		})
	}
}

func AdminUpdateUserHandler(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.PostForm("id")
		name := c.PostForm("name")
		email := c.PostForm("email")
		role := c.PostForm("role")
		pass := c.PostForm("new-password")

		err := userService.UpdateUser(id, name, email, role, pass)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error updating user")
			return
		}

		c.Redirect(http.StatusFound, "/admin/dashboard")
	}
}

func AdminDeleteUserHandler(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := userService.DeleteUser(id)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error deleting user")
			return
		}

		c.Redirect(http.StatusFound, "/admin/dashboard")
	}
}

func AdminLogoutHandler(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.DefaultMany(c, "admin_session")
		session.Clear()
		session.Save()

		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		c.Redirect(http.StatusFound, "/admin/login")
	}
}

func AdminCreateUserHandler(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")
		password := c.PostForm("password")

		err := userService.CreateUser(name, email, password)
		if err != nil {
			c.String(500, "Failed to create user")
			return
		}

		c.Redirect(http.StatusFound, "/admin/dashboard")
	}
}
