package main

import (
	"log"
	"net/http"
	"os"

	"full-domain/internal/database"
	"full-domain/internal/handlers"
	"full-domain/internal/lumberjack"
	"full-domain/internal/middleware"
	"full-domain/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("failed to create log file: " + err.Error())
	}
	defer logFile.Close()

	lumberjack.Init(logFile)
	lumberjack.Logger.Info("starting application")

	lumberjack.Logger.Info("Loading env file")
	if err := godotenv.Load(); err != nil {
		lumberjack.Logger.Error("failed to load env", "error", err)
	}

	if err := database.Connect(); err != nil {
		lumberjack.Logger.Error("database connection failed", "error", err)
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

		session := sessions.DefaultMany(c, "user_session")
		if session.Get("email") != nil && session.Get("role") == "user" {
			c.Redirect(http.StatusFound, "/home")
			return
		}

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
		admin.POST("/update", handlers.AdminUpdateUserHandler(userService))
		admin.GET("/delete/:id", handlers.AdminDeleteUserHandler(userService))
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
