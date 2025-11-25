package main

import (
	"log"
	"net/http"

	"full-domain/internal/database"
	"full-domain/internal/handlers"
	"full-domain/internal/middleware"
	"full-domain/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	userRepo := database.NewUserRepository(database.DB)
	userService := services.NewUserService(userRepo)

	r := gin.Default()

	r.Static("/static", "./static")
	r.LoadHTMLGlob("./template/*")

	store := cookie.NewStore([]byte("super-secret-key"))
	store.Options(sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   3600,
	})

	r.Use(sessions.SessionsMany([]string{
		"user_session",
		"admin_session",
	}, store))

	r.Use(middleware.CacheClear())

	r.GET("/", func(c *gin.Context) {

		session := sessions.DefaultMany(c, "user_session")
		if session.Get("email") != nil && session.Get("role") == "user" {
			c.Redirect(http.StatusFound, "/home")
			return
		}

		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	})

	auth := r.Group("/")
	auth.Use(middleware.AuthRequired(userService), middleware.CacheClear())
	{
		auth.GET("/home", func(c *gin.Context) {

			session := sessions.DefaultMany(c, "user_session")

			name := session.Get("name")
			email := session.Get("email")

			c.HTML(http.StatusOK, "home.tmpl", gin.H{
				"name":  name,
				"email": email,
			})
		})
	}

	r.GET("/admin/login", func(c *gin.Context) {
		session := sessions.DefaultMany(c, "admin_session")
		if session.Get("email") != nil && session.Get("role") == "admin" {
			c.Redirect(http.StatusFound, "/admin/dashboard")
			return
		}
		c.HTML(http.StatusOK, "admin-login.html", nil)
	})

	r.POST("/api/admin/login", handlers.AdminLoginHandler(userService))

	admin := r.Group("/admin")

	admin.Use(middleware.AdminRequired(), middleware.CacheClear())
	{
		admin.GET("/dashboard", handlers.AdminDashboardHandler(userService))
		admin.PATCH("/update", handlers.AdminUpdateUserHandler(userService))
		admin.DELETE("/delete/:id", handlers.AdminDeleteUserHandler(userService))
		admin.POST("/logout", handlers.AdminLogoutHandler(userService))
		admin.POST("/create", handlers.AdminCreateUserHandler(userService))
	}

	api := r.Group("/api")
	{
		api.POST("/login", handlers.LoginHandler(userService))
		api.POST("/signup", handlers.SignupHandler(userService))
		api.POST("/logout", handlers.LogoutHandler(userService))

		user := api.Group("/users")
		user.Use(middleware.AuthRequired(userService))
		{
			user.POST("/password", handlers.UpdatePasswordHandler(userService))
		}
	}

	log.Fatal(r.Run(":8080"))
}
