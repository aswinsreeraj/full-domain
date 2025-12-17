package http

import (
	"full-domain/internal/delivery/http/handlers"
	"full-domain/internal/delivery/http/middleware"
	"full-domain/internal/domain"

	stdhttp "net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func NewRouter(userService domain.UserService) *gin.Engine {
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

	r.Use(middleware.RequestID())
	r.Use(middleware.CacheClear())

	r.GET("/", func(c *gin.Context) {

		session := sessions.DefaultMany(c, "user_session")
		if session.Get("email") != nil && session.Get("role") == "user" {
			c.Redirect(stdhttp.StatusFound, "/home")
			return
		}

		c.HTML(stdhttp.StatusOK, "index.html", nil)
	})

	r.GET("/signup", func(c *gin.Context) {

		session := sessions.DefaultMany(c, "user_session")
		if session.Get("email") != nil && session.Get("role") == "user" {
			c.Redirect(stdhttp.StatusFound, "/home")
			return
		}

		c.HTML(stdhttp.StatusOK, "signup.html", nil)
	})

	auth := r.Group("/")
	auth.Use(middleware.AuthRequired(userService), middleware.CacheClear())
	{
		auth.GET("/home", func(c *gin.Context) {

			session := sessions.DefaultMany(c, "user_session")

			name := session.Get("name")
			email := session.Get("email")

			c.HTML(stdhttp.StatusOK, "home.tmpl", gin.H{
				"name":  name,
				"email": email,
			})
		})
	}

	r.GET("/admin/login", func(c *gin.Context) {
		session := sessions.DefaultMany(c, "admin_session")
		if session.Get("email") != nil && session.Get("role") == "admin" {
			c.Redirect(stdhttp.StatusFound, "/admin/dashboard")
			return
		}
		c.HTML(stdhttp.StatusOK, "admin-login.html", nil)
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

	return r
}
