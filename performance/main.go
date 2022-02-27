package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func worker(wg *sync.WaitGroup){
	defer wg.Done()
	for i :=0;i<100;i++{
		_, err := http.Get("http://localhost:6985/awesome-page")
		if err!=nil{
			log.Println(err)
			continue
		}
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(5)))
	}

}



func main(){
	wg := new(sync.WaitGroup)
	wg.Add(10)
	for i:=0; i<10;i++{
		go worker(wg)
	}
	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm,os.Interrupt)
	go func(){
		wg.Wait()
	}()

	select{
	case <- sigTerm:
		log.Println("SigTerm, exiting !")
		return
	}





}
