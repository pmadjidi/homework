package main

import "github.com/gorilla/mux"

type leaderboard map[string]int
type groups map[string]map[string]bool

const EMPTYSTRING = ""
const NOTFOUND = -1
const TIMEOUT = 2

// potentially 100.000 open http requests, OBS buffered channels
// Obs buffered channels create risk for data loss on server crash...
// Tune for low latency on external http requests

const MAXQUEUELENGTH = 100000

const MAXITERATIONLIMIT = 1000 // concurrent request to API server

// always good to put a bound on datastructures...

const MAXNUMBEROFSTEPSINPUT = 1000
const MAXNUMBERSOFWALKERS = 1000000
const MAXNUMBEROFGROUPS = 100000
const MAXNUMBEROFWALKERSINGROUP = 2000

type pedometers struct {
	name string
	leaderboard
	groups
	leaderBoardCmd         chan *request
	leaderBoardCmdInternal chan *request // internal loop
	groupsCmd              chan *request
	groupsCmdInternal      chan *request // internal loop
}

type App struct {
	*pedometers
	quit chan bool
	*mux.Router
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
