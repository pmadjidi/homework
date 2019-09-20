package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)


func TestReadConfigurationOnDefaults(t *testing.T) {

	os.Setenv("MAXQUEUELENGTH","")
	os.Setenv("MAXITERATIONLIMIT","")
	os.Setenv("MAXNUMBEROFSTEPSINPUT","")
	os.Setenv("MAXNUMBERSOFWALKERS","")
	os.Setenv("MAXNUMBEROFGROUPS","")
	os.Setenv("MAXNUMBEROFWALKERSINGROUP","")
	os.Setenv("TIMEOUT","")

	config := readConfig()

	assert.Equal(t, config.MAXQUEUELENGTH,1 , "shoud be 100000")
	assert.Equal(t, config.MAXITERATIONLIMIT,1000 , "shoud be 1000")
	assert.Equal(t, config.MAXNUMBEROFSTEPSINPUT,1000 , "shoud be 100000")
	assert.Equal(t, config.MAXNUMBERSOFWALKERS,1000000 , "shoud be 100000")
	assert.Equal(t, config.MAXNUMBEROFGROUPS,100000 , "shoud be 100000")
	assert.Equal(t, config.MAXNUMBEROFWALKERSINGROUP,2000 , "shoud be 100000")
	assert.Equal(t, config.TIMEOUT,2 , "shoud be 2")
}

func TestReadConfigurationOnEnv(t *testing.T) {

	os.Setenv("MAXQUEUELENGTH","30000")
	os.Setenv("MAXITERATIONLIMIT","40000")
	os.Setenv("MAXNUMBEROFSTEPSINPUT","50000")
	os.Setenv("MAXNUMBERSOFWALKERS","60000")
	os.Setenv("MAXNUMBEROFGROUPS","70000")
	os.Setenv("MAXNUMBEROFWALKERSINGROUP","80000")
	os.Setenv("TIMEOUT","3")

	config := readConfig()

	assert.Equal(t, config.MAXQUEUELENGTH,30000 , "shoud be 30000")
	assert.Equal(t, config.MAXITERATIONLIMIT,40000 , "shoud be 1000")
	assert.Equal(t, config.MAXNUMBEROFSTEPSINPUT,50000 , "shoud be 40000")
	assert.Equal(t, config.MAXNUMBERSOFWALKERS,60000 , "shoud be 60000")
	assert.Equal(t, config.MAXNUMBEROFGROUPS,70000 , "shoud be 70000")
	assert.Equal(t, config.MAXNUMBEROFWALKERSINGROUP,80000 , "shoud be 80000")
	assert.Equal(t, config.TIMEOUT,3 , "shoud be 3")
}


func TestReadConfigurationOneMissing(t *testing.T) {

	os.Setenv("MAXQUEUELENGTH","30000")
	os.Setenv("MAXITERATIONLIMIT","40000")
	os.Setenv("MAXNUMBEROFSTEPSINPUT","")
	os.Setenv("MAXNUMBERSOFWALKERS","60000")
	os.Setenv("MAXNUMBEROFGROUPS","70000")
	os.Setenv("MAXNUMBEROFWALKERSINGROUP","80000")
	os.Setenv("TIMEOUT","4")

	config := readConfig()

	assert.Equal(t, config.MAXQUEUELENGTH,30000 , "shoud be 30000")
	assert.Equal(t, config.MAXITERATIONLIMIT,40000 , "shoud be 1000")
	assert.Equal(t, config.MAXNUMBEROFSTEPSINPUT,1000 , "shoud be 1000")
	assert.Equal(t, config.MAXNUMBERSOFWALKERS,60000 , "shoud be 60000")
	assert.Equal(t, config.MAXNUMBEROFGROUPS,70000 , "shoud be 70000")
	assert.Equal(t, config.MAXNUMBEROFWALKERSINGROUP,80000 , "shoud be 80000")
	assert.Equal(t, config.TIMEOUT,4 , "shoud be 4")
}



