package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"time"
)

func newApp(name string) *App {
	config := readConfig()
	return &App{
		make(map[int]*pedometers),
		make(chan bool),
		mux.NewRouter(),
		make(chan *request),
		config,
	}
}

func (a *App) startShardHandler(quit chan bool) {
	for i := 0; i < a.config.NUMBEROFSHARDS; i++ {
		a.shards[i] = newPedometers(i, a.config)
		a.shards[i].startPedometers(quit)
	}
}

func (a *App) execLeadBoardCmd(req *request) {
	waitDuration := time.Duration(a.config.TIMEOUT)
	a.steperHash(req)
	index := hextoint(req.Hash[0:SHARTSLICE])


	if (req.Source == EXTERNAL) {
		select {
		case a.shards[index].leaderBoardCmd <- req:
		case <-time.After(waitDuration * time.Second):
			req.Error = &TimeOutError{}
			req.resp <- req
			close(req.resp)
		}
	} else {
		if req.shard == nil {
			select {
			case a.shards[index].leaderBoardCmdInternal <- req:
			case <-time.After(waitDuration * time.Second):
				req.Error = &TimeOutError{}
				req.resp <- req
				close(req.resp)
			}
		} else {
			select {
			case req.shard.leaderBoardCmdInternal <- req:
			case <-time.After(waitDuration * time.Second):
				req.Error = &TimeOutError{}
				req.resp <- req
				close(req.resp)
			}
		}
	}
}

func (a *App) execGroupCmd(req *request) {
	waitDuration := time.Duration(a.config.TIMEOUT)
	a.groupHash(req)
	index := hextoint(req.Hash[0:SHARTSLICE])
	if (req.Source == EXTERNAL) {
		select {
		case a.shards[index].groupsCmd <- req:
		case <-time.After(waitDuration * time.Second):
			req.Error = &TimeOutError{}
			req.resp <- req
			close(req.resp)
		}
	} else {
		if req.shard == nil {
			select {
			case a.shards[index].groupsCmdInternal <- req:
			case <-time.After(waitDuration * time.Second):
				req.Error = &TimeOutError{}
				req.resp <- req
				close(req.resp)
			}
		} else {
			select {
			case req.shard.groupsCmdInternal <- req:
			case <-time.After(waitDuration * time.Second):
				req.Error = &TimeOutError{}
				req.resp <- req
				close(req.resp)
			}
		}
	}
}

func (a *App) start() {
	println("Starting with configuration:")
	PrettyPrint(a.config)
	a.configureRoutes()
	a.startShardHandler(APP.quit)
	a.startWebServer()
}

func (a *App) shutdown() {
	fmt.Println("Exiting...")
	close(a.quit)
}
