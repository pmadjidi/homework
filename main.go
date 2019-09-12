package main

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"time"
)

var APP *App

func shotdown() {
	fmt.Println("Exiting...")
	close(APP.quit)
}


func main() {

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		shotdown()
		<-time.After(TIMEOUT * time.Second)
		os.Exit(1)
	}()

	APP = newApp("Apsis Homework")
	APP.start()

}


