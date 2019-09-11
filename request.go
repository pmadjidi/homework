package main

import "encoding/json"
import "fmt"


type request struct {
	Cmd command  `json:"cmd"`
	Name string `json:"name"`
	Group string `json:"group"`
	Steps int `json:"steps"`
	Error error `json:"error"`
	Result leaderboard
	resp chan *request
}

func newRequest() *request {
	return &request{
		NOP,
		"",
		"",
		0,
		nil,
		nil,
		make(chan *request,1),
	}
}

func (r *request) print() {
	js, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", js)
}

func (r *request) err () {
	fmt.Printf("%s",r.Error.Error())
}
