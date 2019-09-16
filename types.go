package main

import (
	"github.com/gorilla/mux"
)

type leaderboard map[string]int
type groups map[string]map[string]bool

type pedometers struct {
	name string
	leaderboard
	groups
	leaderBoardCmd         chan *request
	leaderBoardCmdInternal chan *request // internal loop
	groupsCmd              chan *request
	groupsCmdInternal      chan *request // internal loop
	config                 *config
}

type config struct {
	MAXQUEUELENGTH    int
	MAXITERATIONLIMIT int // concurrent request to API server
	// always good to put a bound on datastructures...
	MAXNUMBEROFSTEPSINPUT     int
	MAXNUMBERSOFWALKERS       int
	MAXNUMBEROFGROUPS         int
	MAXNUMBEROFWALKERSINGROUP int
	TIMEOUT int
}

type App struct {
	*pedometers
	quit chan bool
	*mux.Router
	*config
}

type command int

type outputStep struct {
	Name  string `json:"name"`
	Steps int    `json:"Steps"`
}

type outputGroup struct {
	Name    string      `json:"name"`
	Steps   int         `json:"Steps"`
	Members leaderboard `json:"members"`
}

type source int
