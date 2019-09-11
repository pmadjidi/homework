package main

import "github.com/gorilla/mux"


func newApp (name string) *App {
	return &App{
		newPedometers(name),
		make(chan bool),
		mux.NewRouter(),
	}
}