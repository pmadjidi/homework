package main

import (
	"github.com/gorilla/mux"
	"fmt"
)


func newApp (name string) *App {
	return &App{
		newPedometers(name),
		make(chan bool),
		mux.NewRouter(),
	}
}

func (a *App) start() {
	a.configureRoutes()
	a.startPedometers(APP.quit)
	a.startWebServer()
}

func (a *App) shotdown() {
	fmt.Println("Exiting...")
	close(a.quit)
}
