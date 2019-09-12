package main

import "github.com/gorilla/mux"

type leaderboard map[string]int
type groups map[string]map[string]bool

const EMPTYSTRING = ""
const NOTFOUND = -1
const TIMEOUT = 2


const MAXQUEUELENGTH = 10000 // works best if 20% more then MAXITERATIONLIMIT
const MAXITERATIONLIMIT = 8000 // concurrent request to API server

type pedometers struct {
	name string
	leaderboard
	groups
	leaderBoardCmd chan *request
	groupsCmd      chan *request
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
	Name  string `json:"name"`
	Steps int    `json:"Steps"`
	Members leaderboard  `json:"members"`
}
