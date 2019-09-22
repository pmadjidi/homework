package main

import "sync"

func newPedometers(name string, config *config) *pedometers {
	var g sync.Map
	var l sync.Map
	return &pedometers{name,l,g,config}
}
