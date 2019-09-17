package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"time"
)

func newApp(name string) *App {
	config := readConfig()
	return &App{
		make(map[string]*pedometers),
		make(chan bool),
		mux.NewRouter(),
		make(chan *request),
		config,
	}
}

func (a *App) shardHandler(quit chan bool) {
	go func() {
		println("Starting Shard Handler")
		for {
			select {
			case req := <- a.Cmd:
				shardKey := req.Hash[0:SHARTSLICE]
				pedomter, ok := a.shards[shardKey]
				if !ok {
					a.shards[shardKey] = newPedometers(shardKey,a.config)
					a.shards[shardKey].startPedometers(quit)
				}
				
			case <-quit:
				println("Stoping group processor")
				return
			default:
			}
		}
	}()
}


func (a *App) execLeadBoardCmd(req *request) {
	waitDuration := time.Duration(a.config.TIMEOUT)
	a.steperHash(req)
	if (req.Source == EXTERNAL) {
		select {
		case p.leaderBoardCmd <- req:
		case <-time.After(waitDuration * time.Second):
			req.Error = &TimeOutError{}
			req.resp <- req
			close(req.resp)
		}
	} else {
		select {
		case p.leaderBoardCmdInternal <- req:
		case <-time.After(waitDuration  * time.Second):
			req.Error = &TimeOutError{}
			req.resp <- req
			close(req.resp)
		}
	}
}

func (a *App) execGroupCmd(req *request) {
	waitDuration := time.Duration(p.config.TIMEOUT)
	a.groupHash(req)
	if (req.Source == EXTERNAL) {
		select {
		case p.groupsCmd <- req:
		case <-time.After(waitDuration * time.Second):
			req.Error = &TimeOutError{}
			req.resp <- req
			close(req.resp)
		}
	} else {
		select {
		case p.groupsCmdInternal <- req:
		case <-time.After(waitDuration  * time.Second):
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
	a.startPedometers(APP.quit)
	a.startWebServer()
}

func (a *App) shutdown() {
	fmt.Println("Exiting...")
	close(a.quit)
}
