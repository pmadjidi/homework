package main

import (
	"fmt"
	"github.com/gorilla/mux"
)

func newApp(name string) *App {
	config := readConfig()
	return &App{
		newPedometers(name, config),
		mux.NewRouter(),
		config,
		make(chan bool),
	}
}

func (a *App) start() {
	println("Starting with configuration:")
	PrettyPrint(a.config)
	a.configureRoutes()
	a.startWebServer()
}

func (a *App) shutdown() {
	fmt.Println("Exiting...")
	close(a.quit)
}
