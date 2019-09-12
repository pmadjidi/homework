package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"runtime"
)

var APP *App



func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		APP.shotdown()
		<-time.After(2 * time.Second)
		os.Exit(1)
	}()

	APP = newApp("Apsis Homework")
	APP.start()

}

