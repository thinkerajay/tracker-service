package config

import (
	"github.com/gin-gonic/gin"
	"github.com/thinkerajay/tracker-service/handlers"
)


func ConfigureRouter(router *gin.Engine){
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/awesome-page/views", handlers.ViewsHandler)
	router.GET("/api/v1/awesome-page/views", handlers.ViewsApiHandler)

	router.GET("/", handlers.HomePageHandler)
	router.GET("/awesome-page",handlers.AwesomePageHandler)
}
