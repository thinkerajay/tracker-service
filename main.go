package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/thinkerajay/tracker-service/config"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT must be set")
	}
	router := gin.New()
	config.ConfigureRouter(router)
	go func (){
		log.Fatalln(router.Run(":" + port))
	}()
	gracefulShutdown()
}
