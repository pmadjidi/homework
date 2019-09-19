package main

import (
	"encoding/json"
	"flag"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"testing"
	"time"
)


func TestMain(m *testing.M) {
	// Pretend to open our DB connection


	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		APP.shutdown()
		<-time.After(2 * time.Second)
		os.Exit(1)
	}()

	APP = newApp("Apsis Homework")
	go func() {
		APP.start()
	}()

	exitCode := m.Run()
	os.Exit(exitCode)
}


func TestRestAPI(t *testing.T) {

	resp, err := http.Get("http://localhost:8080/add/step/payam")
	if err == nil {
		defer resp.Body.Close()
	}

	if err != nil {
		t.Fail()
	} else {
		var result map[string]int
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, result["payam"], 0)
	}

	resp1, err := http.Get("http://localhost:8080/add/group/t1")
	if err == nil {
		defer resp1.Body.Close()
	}


	if err != nil {
		t.Fail()
	} else {
		var result outputGroup
		json.NewDecoder(resp1.Body).Decode(&result)
		assert.Equal(t, result.Name, "t1")
		assert.Equal(t, result.Steps, 0)
	}

	resp2, err := http.Get("http://localhost:8080/extend/t1/payam")
	if err == nil {
		defer resp2.Body.Close()
	}

	if err != nil {
		t.Fail()
	} else {
		var result outputStep
		json.NewDecoder(resp2.Body).Decode(&result)
		assert.Equal(t, result.Name, "payam")
		assert.Equal(t, result.Steps, 0)
	}

	resp3, err := http.Get("http://localhost:8080/inc/payam/10")
	if err == nil {
		defer resp3.Body.Close()
	}


	if err != nil {
		t.Fail()
	} else {
		var result outputStep
		json.NewDecoder(resp3.Body).Decode(&result)
		assert.Equal(t, result.Name, "payam")
		assert.Equal(t, result.Steps, 10)
	}

	resp4, err := http.Get("http://localhost:8080/get/group/t1")
	if err == nil {
		defer resp4.Body.Close()
	}

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
	if err == nil {
		defer resp5.Body.Close()
	}

	if err != nil {
		t.Fail()
	} else {
		var result outputGroup
		json.NewDecoder(resp5.Body).Decode(&result)
		assert.Equal(t, result.Name, "t2")
		assert.Equal(t, result.Steps, 0)
	}


	resp6, err := http.Get("http://localhost:8080/get/allgroups")
	if err == nil {
		defer resp5.Body.Close()
	}

	if err != nil {
		t.Fail()
	} else {
		var result outputGroup
		json.NewDecoder(resp6.Body).Decode(&result)
		assert.Equal(t, result.Name, "t2")
		assert.Equal(t, result.Steps, 0)
	}


	//and so furth....

}
