package main

import (
	"github.com/gorilla/mux"
	"sync"
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
	SHARDS int
	HASHBITSTOSHARD int
	PORT int
}

type shards map[int]*pedometers

type App struct {
	shards
	begin chan bool
	quit chan bool
	*mux.Router
	Cmd      chan *request
	*config
	cache sync.Map
}

type command int

type outputStep struct {
	Name  string `json:"name"`
	Points int    `json:"Points"`
}

type outputGroup struct {
	Name    string      `json:"name"`
	Points   int         `json:"Points"`
	Members leaderboard `json:"members"`
}

type source int

type request struct {
	Cmd     command `json:"cmd"`
	Source  source  `json:"source"`
	Name    string  `json:"name"`
	Group   string  `json:"group"`
	Points   int     `json:"points"`
	Error   error   `json:"error"`
	Hash    string  `json:"hash"`
	Result  leaderboard
	Results map[string]leaderboard
	resp    chan *request
	index   int
}
