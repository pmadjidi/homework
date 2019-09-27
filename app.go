package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"sync"
	"time"
)

func newApp(name string) *App {
	config := readConfig()
	var cache sync.Map
	return &App{
		make(map[int]*pedometers),
		make(chan bool),
		make(chan bool),
		mux.NewRouter(),
		make(chan *request),
		config,
		cache,
	}
}

func (a *App) startShardHandler(quit chan bool) {

	for i := 0; i < a.config.SHARDS ; i++ {
		a.shards[i] = newPedometers(i, a.config)
		a.shards[i].startPedometers(quit)
	}
	go func() {
		a.begin <-true
	}()

}

func (a *App) execLeadBoardCmd(req *request) {
	var index int
	waitDuration := time.Duration(a.config.TIMEOUT)

	if req.index >= 0 {
		index = req.index
	} else {
		a.userHash(req)
		index = hextoint(req.Hash[0:a.config.HASHBITSTOSHARD])
		req.index = index
		println(req.Cmd.String()," To Node: ",req.index)
	}

	if (req.Source == EXTERNAL) {
		select {
		case a.shards[index].leaderBoardCmd <- req:
		case <-time.After(waitDuration * time.Second):
			req.Error = &TimeOutError{}
			req.resp <- req
			close(req.resp)
		}
	} else {
		select {
		case a.shards[index].leaderBoardCmdInternal <- req:
		case <-time.After(waitDuration * time.Second):
			req.Error = &TimeOutError{}
			req.resp <- req
			close(req.resp)
		}
	}
}

func (a *App) execGroupCmd(req *request) {
	var index int

	if req.index >= 0 {
		index = req.index
	} else {
		a.groupHash(req)
		index = hextoint(req.Hash[0:a.config.HASHBITSTOSHARD])
		req.index = index
		println(req.Cmd.String()," To Node: ",req.index )
	}

	if (req.Source == EXTERNAL) {
		select {
		case a.shards[index].groupsCmd <- req:
		case <-time.After(time.Duration(a.config.TIMEOUT) * time.Second):
			req.Error = &TimeOutError{}
			req.resp <- req
			close(req.resp)
		}
	} else {
		select {
		case a.shards[index].groupsCmdInternal <- req:
		case <-time.After(time.Duration(a.config.TIMEOUT) * time.Second):
			req.Error = &TimeOutError{}
			req.resp <- req
			close(req.resp)
		}
	}
}

func (a *App) start() {
	println("Starting with configuration:")
	PrettyPrint(a.config)
	a.configureRoutes()
	a.startShardHandler(APP.quit)
	a.startWebServer()
	<- a.begin
}

func (a *App) shutdown() {
	fmt.Println("Exiting...")
	close(a.quit)
}
