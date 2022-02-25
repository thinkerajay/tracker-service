package config

import (
	"github.com/gin-gonic/gin"
	"github.com/thinkerajay/tracker-service/handlers"
)


func ConfigureRouter(router *gin.Engine){
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/views", handlers.ViewsHandler)

	router.GET("/", handlers.HomePageHandler)
	router.GET("/awesome-page",handlers.AwesomePageHandler)
}
