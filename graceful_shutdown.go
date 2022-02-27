package main

import (
	"context"
	"github.com/thinkerajay/tracker-service/db_service"
	"log"
	"os"
	"os/signal"
)

func gracefulShutdown() {
	//graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	go db_service.ConsumePageViewEvents(ctx)

	select {
	case <-signals:
		cancel()
		log.Println("SigTerm, exiting the application !")
	}
}
