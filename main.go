package main

import (
	"context"
	"github.com/thinkerajay/tracker-service/config"
	"github.com/thinkerajay/tracker-service/db_service"
	"log"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT must be set")
	}
	port = "6985"
	router := gin.New()
	config.ConfigureRouter(router)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	go db_service.Consume(ctx)
	go func (){
		log.Fatalln(router.Run(":" + port))
	}()
	select{
	case <- signals:
		cancel()
		log.Println("SigTerm, exiting the application !")
	}

}
