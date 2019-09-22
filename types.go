package main

import (
	"github.com/gorilla/mux"
	"sync"
)



type pedometers struct {
	name string
	leaderboard sync.Map
	groups sync.Map
	config  *config
}

type config struct {
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
	*mux.Router
	*config
}

type command int

type outputStep struct {
	Name  string `json:"name"`
	Steps int    `json:"Steps"`
}
type outputGroup struct {
	Name  string `json:"name"`
	Group string    `json:"group"`
}



type outputGroupMembers struct {
	Name    string      `json:"name"`
	Steps   int         `json:"Steps"`
	Members map[string]int `json:"members"`
}

type source int
