package main

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
)

func TestAddWalkerFail(t *testing.T) {
	e := APP.AddWalker("")
	assert.Equal(t, e.Error(), "NO_NAME", "Should generate NO_NAME")
}

func TestAddWalkerSucess(t *testing.T) {

	e := APP.AddWalker("Payam")
	assert.Equal(t, e, nil, "No Error")

	result := APP.ListAllSteppers()
	assert.Equal(t, result["Payam"], 0, "Counter set to zer0")
}

func TestAddWalker(t *testing.T) {

	e := APP.AddWalker("Payam")
	assert.Equal(t,e, nil, "No Error")


	APP.AddWalker("Mikael")
	result := APP.ListAllSteppers()
	assert.Equal(t, result["Payam"], 0, "Counter set to zer0")
	assert.Equal(t, result["Mikael"], 0, "Counter set to zer0")
}

func TestAddWalkerFailExits(t *testing.T) {
	name := "Payam"
	e := APP.AddWalker(name)
	assert.Equal(t, e, nil, "No Error")


	e := APP.AddWalker(name)
	assert.Equal(t,e.Error(), "NAME_EXISTS", "Should generate NAME_EXISTS")
}

func TestDeleteWalkerMissing(t *testing.T) {
	e := APP.DeleteWalker("Payam")
	assert.EqualError(t,e, "COMMAND_NOT_IMPLEMENTED", "Should generate COMMAND_NOT_IMPLEMENTED")
}

func TestConcurrentWalker(t *testing.T) {
	name := "Payam"
	APP.AddWalker(name)
	var waitgroup sync.WaitGroup
	for i := 0; i < APP.config.MAXITERATIONLIMIT; i++ {
		waitgroup.Add(1)
		go func() {
			APP.RegisterSteps(name,1)
			waitgroup.Done()
		}()
	}
	waitgroup.Wait()
	res := APP.ListAllSteppers()
	assert.Equal(t, result["Payam"], APP.config.MAXITERATIONLIMIT, "Payam took all the steps")

}

func TestConcurrentWalkerWithStepArgument(t *testing.T) {
	name := "Payam"
	APP.AddWalker(name)
	var waitgroup sync.WaitGroup
	for i := 0; i < APP.config.MAXITERATIONLIMIT; i++ {
		waitgroup.Add(1)
		go func() {
			APP.RegisterSteps(name,5)
			waitgroup.Done()
		}()
	}
	waitgroup.Wait()
	result := APP.ListAllSteppers()

	assert.Equal(t, result["Payam"], APP.config.MAXITERATIONLIMIT*5, "Payam took all the steps")

}

func TestConcurrentWalkerWithStepArgukmentOne(t *testing.T) {
	var waitgroup sync.WaitGroup
	for i := 0; i < APP.config.MAXITERATIONLIMIT; i++ {
		waitgroup.Add(1)
		go func(index int) {
			name := "Payam" + strconv.Itoa(index)
			APP.AddWalker(name)
			APP.RegisterSteps(name,10)
			waitgroup.Done()
		}(i)
	}
	waitgroup.Wait()

	result := APP.ListAllSteppers()

	for i := 0; i < APP.config.MAXITERATIONLIMIT; i++ {
		assert.Equal(t, result["Payam"+strconv.Itoa(i)], 10, "Each walker takes only 10 steps")
	}


}

func TestConcurrentWalkerAPIrace(t *testing.T) {
	var waitgroup sync.WaitGroup
	for i := 0; i < APP.config.MAXITERATIONLIMIT; i++ {
		waitgroup.Add(1)
		go func(index int) {
			name  := "Payam" + strconv.Itoa(index)
			APP.AddWalker(name)
			waitgroup.Done()
		}(i)
	}

	waitgroup.Wait()
	for i := 0; i < APP.config.MAXITERATIONLIMIT; i++ {
		waitgroup.Add(1)
		go func(index int) {
			name := "Payam" + strconv.Itoa(index)
			APP.RegisterSteps(name ,10)
			waitgroup.Done()
		}(i)

		waitgroup.Add(1)
		go func(index int) {
			name := "Payam" + strconv.Itoa(index)
			APP.RegisterSteps(name,10)
			waitgroup.Done()
		}(i)

	}

	waitgroup.Wait()
	println("callculating resutls")

	result := APP.ListAllSteppers()

	var missed = 0

	for i := 0; i < APP.config.MAXITERATIONLIMIT; i++ {
		if result["Payam"+strconv.Itoa(i)] != 20 {
			missed++
		}
	}

	println("***********************************************", missed)

	assert.Equal(t, missed, 0, "Can not increment stepcounter before creation")
	
}

func TestAddGrouprFail(t *testing.T) {

	e := APP.AddGroup("")
	assert.EqualError(t, e, "NO_GROUP", "Should generate NO_GROUP")
}

func TestAddGrouprSuccess(t *testing.T) {

	Group := "Apsis"
	e := APP.AddGroup(Group)
	assert.Equal(t, e, nil, "Should generate NO_GROUP")
}

func TestAddGrouprExistFail(t *testing.T) {



	Group := "Apsis"
	APP.AddGroup(Group)

	e := APP.AddGroup(Group)


	assert.EqualError(t, e, "GROUP_EXISTS", "Should generate GROUP_EXISTS")
}

func TestAddWalkerToGroupSucess(t *testing.T) {
	e := APP.AddWalker("Payam")
	assert.Equal(t, e, nil, "No Error")
	e = APP.AddGroup("KTH")
	assert.Equal(t, e, nil, "No Error")


	e = APP.AddWalkerToGroup("Payam","KTH")
	assert.Equal(t,e , nil, "No Error")


	result := APP.ListAllSteppers()
	assert.Equal(t, result["Payam"], 0, "Counter set to zer0")
}

func TestAddWalkerToGroupFail(t *testing.T) {



	e := APP.AddWalker("Payam")
	assert.Equal(t, e, nil, "No Error")



	APP.AddGroup("Apsis")

	APP.AddWalkerToGroup("Payam","KTH")


	assert.EqualError(t, e, "GROUP_MISSING", "Should generate GROUP_MISSING")

}

func TestAddMultpleWalkersToMultipleGroups(t *testing.T) {
	var waitgroup sync.WaitGroup
	APP.AddGroup("TestGroup")
	for i := 0; i < APP.config.MAXITERATIONLIMIT; i++ {
		waitgroup.Add(1)
		go func(index int) {

			name := "Payam" + strconv.Itoa(index)
			APP.AddWalker(name)

			APP.RegisterSteps(name,10)
			APP.AddWalkerToGroup(name,"TestGroup")
			waitgroup.Done()
		}(i)
	}

	waitgroup.Wait()


	result, e := APP.ListGroup("TestGroup")


	assert.Equal(t, res.Steps, p.config.MAXITERATIONLIMIT*10, "Testgroup should add up")

	req = newRequest()
	p.ListAll(req)
	res = <-req.resp
	assert.Equal(t, res.Steps, p.config.MAXITERATIONLIMIT*10, "List all should add up")

}
