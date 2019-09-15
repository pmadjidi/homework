package main

const EMPTYSTRING = ""
const NOTFOUND = -1
const TIMEOUT = 2

// potentially 100.000 open http requests, OBS buffered channels
// Obs buffered channels create risk for data loss on server crash...
// Tune for low latency on external http requests

const MAXQUEUELENGTH = 100000

const MAXITERATIONLIMIT = 1000 // concurrent request to API server

// always good to put a bound on datastructures...

const MAXNUMBEROFSTEPSINPUT = 1000
const MAXNUMBERSOFWALKERS = 1000000
const MAXNUMBEROFGROUPS = 100000
const MAXNUMBEROFWALKERSINGROUP = 2000

