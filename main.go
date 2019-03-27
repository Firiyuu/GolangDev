package main

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"

	"pinehq/controllers"
	"pinehq/models"
	"pinehq/system"
	
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/utrack/gin-csrf"
)

func main() {
	initLogger()
	system.LoadConfig()
	models.SetDB(system.GetConnectionString())
	models.AutoMigrate()
	system.LoadTemplates()



    // SETUP USER ROUTER
	router := gin.Default()
	router.SetHTMLTemplate(system.GetTemplates())

	//setup sessions
	config := system.GetConfig()
	store := memstore.NewStore([]byte(config.SessionSecret))
	store.Options(sessions.Options{HttpOnly: true, MaxAge: 7 * 86400}) //Also set Secure: true if using SSL, you should though
	router.Use(sessions.Sessions("gin-session", store))

	//setup csrf protection
	router.Use(csrf.Middleware(csrf.Options{
		Secret: config.SessionSecret,
		ErrorFunc: func(c *gin.Context) {
			logrus.Error("CSRF token mismatch")
			controllers.ShowErrorPage(c, 400, fmt.Errorf("CSRF token mismatch"))
			c.Abort()
		},
	}))

	router.StaticFS("/public", http.Dir(system.PublicPath()))
	router.Use(controllers.ContextData())

	router.GET("/", controllers.HomeGet)
	router.NoRoute(controllers.NotFound)
	router.NoMethod(controllers.MethodNotAllowed)

	if system.GetConfig().SignupEnabled {
		router.GET("/signup", controllers.SignUpGet)
		router.POST("/signup", controllers.SignUpPost)
	}
	router.GET("/signin", controllers.SignInGet)
	router.POST("/signin", controllers.SignInPost)
	router.GET("/logout", controllers.LogoutGet)


    // AUTHORIZATION FIRST
	authorized := router.Group("/admin")
	authorized.Use(controllers.AuthRequired())
	{   
		authorized.GET("/", controllers.AdminGet)
		authorized.Use(controllers.Mailer())

		// MAILING
		authorized.GET("/notifications", controllers.ShowMail)
	    authorized.POST("/notifications", controllers.SendMail)
	    authorized.GET("/notifications/list", controllers.EmailGet)

		authorized.GET("/notifications/list/:id/edit", controllers.EmailEdit)
		authorized.POST("/notifications/list/:id/edit", controllers.EmailUpdate)
	    



	}

	router.Run(":8080")
}

//initLogger initializes logrus logger with some defaults
func initLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	if gin.Mode() == gin.DebugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
