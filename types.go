package main

import "github.com/gorilla/mux"

type leaderboard map[string]int
type groups map[string]map[string]bool

const EMPTYSTRING = ""
const NOTFOUND = -1
const TIMEOUT = 2

/*

on 8 core macbookpro 2018 with 32GIG memory, 40.000 concurrent connections....

Stoping leaderboard processor
Stoping group processor
PASS
ok  	github.com/pmadjidi/homework	17.186s
 */


const MAXQUEUELENGTH = 100000 // potentially 100.000 open http requests....
const MAXITERATIONLIMIT = 10000// concurrent request to API server

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
