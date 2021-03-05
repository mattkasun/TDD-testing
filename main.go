package main

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func main() {
	router := SetupRouter()
	//Serve the app
	router.Run()
}

func SetupRouter() *gin.Engine {
	gin.DisableConsoleColor()
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	router := gin.Default()

	store := memstore.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("todo", store))

	//process templates
	router.LoadHTMLGlob("html/*")
	//Initialize routes
	initRoutes(router)
	return router
}

func initRoutes(router *gin.Engine) {
	//routes
	router.Static("/stylesheet", "./stylesheet")
	router.StaticFile("favicon.ico", "./resources/favicon.ico")
	router.GET("/", displayLanding)
	router.GET("/register", displayRegister)
	router.GET("/login", displayLogin)
	router.POST("/register", processRegistration)
	router.POST("/login", processLogin)

	private := router.Group("/", authRequired())
	{
		private.GET("/main", displayMainPage)
		private.GET("/add", displayAddPage)
		private.POST("/add", processAddTask)
	}
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		loggedIn := session.Get("loggedIn")
		if loggedIn == nil {
			c.HTML(http.StatusUnauthorized, "login", nil)
			c.Abort()
		}
	}
}
