package main



var APP *App

func main() {

	APP = newApp("Apsis Homework")
	APP.configureRoutes()
	APP.startPedometers(APP.quit)
	APP.startWebServer()
}


