package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"runtime"
	"testing"
)

func TestRestAPI(t *testing.T) {

	runtime.GOMAXPROCS(runtime.NumCPU())
	var quit = make(chan bool)

	go func() {
		APP = newApp("Apsis Homework")
		APP.start()
		<-quit
		APP.shotdown()
	}()

	defer close(quit)

	resp, err := http.Get("http://localhost:8080/add/step/payam")
	defer resp.Body.Close()

	if err != nil {
		t.Fail()
	} else {
		var result map[string]int
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, result["payam"], 0)
	}

	resp1, err := http.Get("http://localhost:8080/add/group/t1")
	defer resp1.Body.Close()

	if err != nil {
		t.Fail()
	} else {
		var result outputGroup
		json.NewDecoder(resp1.Body).Decode(&result)
		assert.Equal(t, result.Name, "t1")
		assert.Equal(t, result.Steps, 0)
	}

	resp2, err := http.Get("http://localhost:8080/extend/t1/payam")
	defer resp2.Body.Close()

	if err != nil {
		t.Fail()
	} else {
		var result outputStep
		json.NewDecoder(resp2.Body).Decode(&result)
		assert.Equal(t, result.Name, "payam")
		assert.Equal(t, result.Steps, 0)
	}

	resp3, err := http.Get("http://localhost:8080/inc/payam/10")
	defer resp3.Body.Close()

	if err != nil {
		t.Fail()
	} else {
		var result outputStep
		json.NewDecoder(resp3.Body).Decode(&result)
		assert.Equal(t, result.Name, "payam")
		assert.Equal(t, result.Steps, 10)
	}

	resp4, err := http.Get("http://localhost:8080/get/group/t1")
	defer resp4.Body.Close()

	if err != nil {
		t.Fail()
	} else {
		var result outputGroup
		json.NewDecoder(resp4.Body).Decode(&result)
		assert.Equal(t, result.Name, "t1")
		assert.Equal(t, result.Steps, 10)
		assert.Equal(t, result.Members["payam"], 10)
	}

	resp5, err := http.Get("http://localhost:8080/add/group/t2")
	defer resp5.Body.Close()

	if err != nil {
		t.Fail()
	} else {
		var result outputGroup
		json.NewDecoder(resp5.Body).Decode(&result)
		assert.Equal(t, result.Name, "t2")
		assert.Equal(t, result.Steps, 0)
	}

	//and so furth....


}
