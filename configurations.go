package main

import (
	"os"
	"strconv"
)

// potentially 100.000 open http requests, OBS buffered channels
// Obs buffered channels create risk for data loss on server crash...
// Tune for low latency on external http requests

//const MAXQUEUELENGTH = 100000

//const MAXITERATIONLIMIT = 1000 // concurrent request to API server

// always good to put a bound on datastructures...

//const MAXNUMBEROFSTEPSINPUT = 1000
//const MAXNUMBERSOFWALKERS = 1000000
//const MAXNUMBEROFGROUPS = 100000
//const MAXNUMBEROFWALKERSINGROUP = 2000

func readIntFromEnv(name string) (int) {
	var result int

	MAXITERATIONLIMIT := 1000

	MAXNUMBEROFSTEPSINPUT := 1000
	MAXNUMBERSOFWALKERS := 1000000
	MAXNUMBEROFGROUPS := 100000
	MAXNUMBEROFWALKERSINGROUP := 2000
	TIMEOUT := 2

	fromEnv := os.Getenv(name)
	iEnv, err := strconv.Atoi(fromEnv)

	switch name {
	case "MAXITERATIONLIMIT":
		if err != nil {
			result = MAXITERATIONLIMIT
		} else {
			result = iEnv
		}
	case "MAXNUMBEROFSTEPSINPUT":
		if err != nil {
			result = MAXNUMBEROFSTEPSINPUT
		} else {
			result = iEnv
		}
	case "MAXNUMBERSOFWALKERS":
		if err != nil {
			result = MAXNUMBERSOFWALKERS
		} else {
			result = iEnv
		}
	case "MAXNUMBEROFGROUPS":
		if err != nil {
			result = MAXNUMBEROFGROUPS
		} else {
			result = iEnv
		}
	case "MAXNUMBEROFWALKERSINGROUP":
		if err != nil {
			result = MAXNUMBEROFWALKERSINGROUP
		} else {
			result = iEnv
		}
	case "TIMEOUT":
		if err != nil {
			result = TIMEOUT
		} else {
			result = iEnv
		}
	}
	return result
}

func readConfig() *config {
	return &config{
		readIntFromEnv("MAXITERATIONLIMIT"),
		readIntFromEnv("MAXNUMBEROFSTEPSINPUT"),
		readIntFromEnv("MAXNUMBERSOFWALKERS"),
		readIntFromEnv("MAXNUMBEROFGROUPS"),
		readIntFromEnv("MAXNUMBEROFWALKERSINGROUP"),
		readIntFromEnv("TIMEOUT"),
	}
}
