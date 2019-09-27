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

	resp, err := http.Get("http://localhost:8090/add/user/payam")
	if err == nil {
		defer resp.Body.Close()
	}

	if err != nil {
		t.Fail()
	} else {
		var result outputStep
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t,result.Name, "payam")
		assert.Equal(t, result.Points, 0)
	}

	resp1, err := http.Get("http://localhost:8090/add/group/t1")
	if err == nil {
		defer resp1.Body.Close()
	}


	if err != nil {
		t.Fail()
	} else {
		var result outputGroup
		json.NewDecoder(resp1.Body).Decode(&result)
		assert.Equal(t, result.Name, "t1")
		assert.Equal(t, result.Points, 0)

	}

	resp2, err := http.Get("http://localhost:8090/extend/t1/payam")
	if err == nil {
		defer resp2.Body.Close()
	}

	if err != nil {
		t.Fail()
	} else {
		var result outputStep
		json.NewDecoder(resp2.Body).Decode(&result)
		assert.Equal(t, result.Name, "payam")
		assert.Equal(t, result.Points, 0)
	}

	resp3, err := http.Get("http://localhost:8090/inc/payam/10")
	if err == nil {
		defer resp3.Body.Close()
	}


	if err != nil {
		t.Fail()
	} else {
		var result outputStep
		json.NewDecoder(resp3.Body).Decode(&result)
		assert.Equal(t, result.Name, "payam")
		assert.Equal(t, result.Points, 10)
	}

	resp4, err := http.Get("http://localhost:8090/get/group/t1")
	if err == nil {
		defer resp4.Body.Close()
	}

	if err != nil {
		t.Fail()
	} else {
		var result outputGroup
		json.NewDecoder(resp4.Body).Decode(&result)
		assert.Equal(t, result.Name, "t1")
		assert.Equal(t, result.Points, 10)
		assert.Equal(t, result.Members["payam"], 10)
	}

	resp5, err := http.Get("http://localhost:8090/add/group/t2")
	if err == nil {
		defer resp5.Body.Close()
	}

	if err != nil {
		t.Fail()
	} else {
		var result outputGroup
		json.NewDecoder(resp5.Body).Decode(&result)
		assert.Equal(t, result.Name, "t2")
		assert.Equal(t, result.Points, 0)
	}


	resp6, err := http.Get("http://localhost:8090/get/groups")
	if err == nil {
		defer resp5.Body.Close()
	}

	if err != nil {
		t.Fail()
	} else {
		var result map[string]leaderboard
		json.NewDecoder(resp6.Body).Decode(&result)
		assert.NotEmpty(t,result)
		PrettyPrint(result)
		t1Sum :=  result["t1"]["SUM"]
		t2Sum := result["t2"]["SUM"]
		println(t1Sum,t2Sum)
		assert.Equal(t,t1Sum, 10)
		assert.Equal(t,t2Sum, 0)
	}


	//and so furth....

}
