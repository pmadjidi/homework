package main

import "github.com/gorilla/mux"

type leaderboard  map[string]int
type groups map[string]map[string]bool
const EMPTYSTRING = ""
const NOTFOUND = -1
const TIMEOUT = 2


type Project struct {
	Id int64 `json:"project_id"`
	Title string `json:"title"`
	Name string `json:"name"`
}


type pedometers struct {
	name string
	leaderboard
	groups
	leaderBoardCmd chan *request
	groupsCmd chan *request
}


type App struct {
	*pedometers
	quit chan bool
	*mux.Router
}


type command int

const (
	NOP command = iota
	ADDWALKER
	REGISTERSTEPS
	GETWALKER
	ADDGROUP
	ADDWALKERTOGROUP
	DELETEWALKER
	RESETSTEPS
	LISTGROUP
	LISTALL
	SCAN
)





