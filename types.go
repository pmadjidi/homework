package main

import (
	"github.com/gorilla/mux"
)

type leaderboard map[string]int
type groups map[string]map[string]int

type pedometers struct {
	index int
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
	NUMBEROFSHARDS int
}

type shards map[int]*pedometers

type App struct {
	shards
	quit chan bool
	*mux.Router
	Cmd      chan *request
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

type request struct {
	Cmd     command `json:"cmd"`
	Source  source  `json:"source"`
	Name    string  `json:"name"`
	Group   string  `json:"group"`
	Steps   int     `json:"steps"`
	Error   error   `json:"error"`
	Hash    string  `json:"hash"`
	Result  leaderboard
	Results map[string]leaderboard
	resp    chan *request
	shard  *pedometers
	index   int
}
