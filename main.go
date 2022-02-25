package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		userId, err := c.Cookie("_tsuid")
		if err != nil {
			c.SetCookie("_tsuid", "random_id", 9000, "/", "herokuapp.com", true, true)
		}
		log.Println(userId)
		c.JSON(200, userId)

	})

	router.Run(":" + port)

}
