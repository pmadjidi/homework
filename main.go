package main

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var APP *App

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		APP.shutdown()
		<-time.After(2 * time.Second)
		os.Exit(1)
	}()
	APP = newApp("Apsis Homework")
	APP.start()
	<-APP.quit

}
