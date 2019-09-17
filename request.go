package main

import "encoding/json"
import "fmt"






func newRequest() *request {
	return &request{
		NOP,
		EXTERNAL,
		"",
		"",
		0,
		nil,
		"",
		nil,
		nil,
		make(chan *request, 1),
		nil,
	}
}

func newRequestInternal() *request {
	return &request{
		NOP,
		INTERNAL,
		"",
		"",
		0,
		nil,
		"",
		nil,
		nil,
		make(chan *request, 1),
		nil,
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

func (r *request) err() {
	fmt.Printf("%s", r.Error.Error())
}
