package main

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"sync"
)


const MAXITERATIONLIMIT = 2000


func TestAddWalkerFail(t *testing.T) {
	quit := make(chan bool)
	p := newPedometers("TestPedometers")
	p.startPedometers(quit)
	defer close(quit)

	req := newRequest()
	go p.AddWalker(req)
	res := <- req.resp
	assert.EqualError(t, res.Error,"NO_NAME","Should generate NO_NAME")
}

func TestAddWalkerSucess(t *testing.T) {

	quit := make(chan bool)
	p := newPedometers("TestPedometers")
	p.startPedometers(quit)
	defer close(quit)

	req := newRequest()
	req.Name = "Payam"
	go p.AddWalker(req)
	res := <- req.resp
	assert.Equal(t, res.Error,nil, "No Error")

	req = newRequest()
	go p.ListAll(req)
	res = <- req.resp
	assert.Equal(t, res.Result["Payam"],0, "Counter set to zer0")
}

func TestAddWalker(t *testing.T) {
	quit := make(chan bool)
	p := newPedometers("TestPedometers")
	p.startPedometers(quit)
	defer close(quit)

	req := newRequest()
	req.Name = "Payam"
	go p.AddWalker(req)
	res := <- req.resp
	assert.Equal(t, res.Error,nil, "No Error")

	req = newRequest()
	req.Name = "Mikael"
	go p.AddWalker(req)
	res = <- req.resp
	assert.Equal(t, res.Result["Payam"],0, "Counter set to zer0")
	assert.Equal(t, res.Result["Mikael"],0, "Counter set to zer0")
	req = newRequest()
	go p.ListAll(req)
	res = <- req.resp
	res.print()
}

func TestAddWalkerFailExits(t *testing.T) {
	quit := make(chan bool)
	p := newPedometers("TestPedometers")
	p.startPedometers(quit)
	defer close(quit)
	req := newRequest()
	req.Name = "Payam"
	go p.AddWalker(req)
	res := <- req.resp
	assert.Equal(t, res.Error,nil, "No Error")
	req = newRequest()
	req.Name = "Payam"
	go p.AddWalker(req)
	res = <- req.resp
	assert.Equal(t, res.Error.Error(),"NAME_EXISTS", "Should generate NAME_EXISTS")
}



func TestDeleteWalkerMissing(t *testing.T) {
	quit := make(chan bool)
	p := newPedometers("TestPedometers")
	p.startPedometers(quit)
	defer close(quit)
	req := newRequest()
	req.Name = "Payam"
	go p.DeleteWalker(req)
	res := <- req.resp
	assert.EqualError(t, res.Error,"COMMAND_NOT_IMPLEMENTED","Should generate COMMAND_NOT_IMPLEMENTED")
	req = newRequest()
	req.Name = "Mikael"
	go p.DeleteWalker(req)
	res = <- req.resp
	assert.EqualError(t, res.Error,"COMMAND_NOT_IMPLEMENTED","Should generate COMMAND_NOT_IMPLEMENTED")

}




func TestConcurrentWalker(t *testing.T) {
	quit := make(chan bool)
	p := newPedometers("TestPedometers")
	p.startPedometers(quit)
	defer close(quit)
	req := newRequest()
	req.Name = "Payam"
	go p.AddWalker(req)
	<- req.resp
	var waitgroup sync.WaitGroup
	for i := 0; i < MAXITERATIONLIMIT; i++ {
		waitgroup.Add(1)
		go func() {
			req := newRequest()
			req.Name = "Payam"
			p.RegisterSteps(req)
			<- req.resp
			waitgroup.Done()
		} ()
	}
	waitgroup.Wait()

	req = newRequest()
	go p.ListAll(req)
	res := <- req.resp
	res.print()

	assert.Equal(t, res.Error,nil, "Expect no Error")
	assert.Equal(t, res.Steps,MAXITERATIONLIMIT, "All the steps taken should add up")
	assert.Equal(t, res.Result["Payam"],MAXITERATIONLIMIT, "Payam took all the steps")

}




func TestConcurrentWalkerWithStepArgument(t *testing.T) {
	quit := make(chan bool)
	p := newPedometers("TestPedometers")
	p.startPedometers(quit)
	defer close(quit)
	req := newRequest()
	req.Name = "Payam"
	go p.AddWalker(req)
	<- req.resp
	var waitgroup sync.WaitGroup
	for i := 0; i < MAXITERATIONLIMIT; i++ {
		waitgroup.Add(1)
		go func() {
			req := newRequest()
			req.Name = "Payam"
			req.Steps = 5
			p.RegisterSteps(req)
			_ = <- req.resp
			waitgroup.Done()
		} ()
	}
	waitgroup.Wait()
	req = newRequest()
	go p.ListAll(req)
	res := <- req.resp


	assert.Equal(t, res.Error,nil, "Expect no Error")
	assert.Equal(t, res.Steps,MAXITERATIONLIMIT * 5, "All the steps taken should add up")
	assert.Equal(t, res.Result["Payam"],MAXITERATIONLIMIT *5, "Payam took all the steps")

}




func TestConcurrentWalkerWithStepArgukmentOne(t *testing.T) {

	quit := make(chan bool)
	p := newPedometers("TestPedometers")
	p.startPedometers(quit)
	defer close(quit)
	req := newRequest()
	var waitgroup sync.WaitGroup
	for i := 0; i < MAXITERATIONLIMIT; i++ {
		waitgroup.Add(1)
		go func(index int) {
			req := newRequest()
			req.Name = "Payam" + strconv.Itoa(index)
			p.AddWalker(req)
			_ = <- req.resp

			req = newRequest()
			req.Name = "Payam" + strconv.Itoa(index)
			req.Steps = 10
			p.RegisterSteps(req)
			_ = <- req.resp
			waitgroup.Done()
		} (i)
	}
	waitgroup.Wait()
	req = newRequest()
	go p.ListAll(req)
	res := <- req.resp

	assert.Equal(t, res.Error,nil, "Expect no Error")
	for i := 0; i < MAXITERATIONLIMIT; i++ {
		assert.Equal(t, res.Result["Payam"  + strconv.Itoa(i)],10, "Each walker takes only 10 steps")
	}

	assert.Equal(t, res.Steps,MAXITERATIONLIMIT * 10, "Each walker takes only 10 steps")


}


func TestConcurrentWalkerAPIrace(t *testing.T) {

	quit := make(chan bool)
	p := newPedometers("TestPedometers")
	p.startPedometers(quit)
	defer close(quit)
	var waitgroup sync.WaitGroup
	for i := 0; i < MAXITERATIONLIMIT; i++ {
		waitgroup.Add(1)
		go func(index int) {
			req := newRequest()
			req.Name = "Payam" + strconv.Itoa(index)
			p.AddWalker(req)
			_ = <-req.resp
			waitgroup.Done()
		}(i)
		waitgroup.Add(1)
		go func(index int) {
			req := newRequest()
			req.Name = "Payam" + strconv.Itoa(index)
			req.Steps = 10
			p.RegisterSteps(req)
			_ = <- req.resp
			waitgroup.Done()
		} (i)

		waitgroup.Add(1)
		go func(index int) {
			req := newRequest()
			req.Name = "Payam" + strconv.Itoa(index)
			req.Steps = 10
			p.RegisterSteps(req)
			_ = <- req.resp
			waitgroup.Done()
		} (i)

	}

	waitgroup.Wait()
	println("callculating resutls")
	all := newRequest()
	go p.ListAll(all)
	res := <- all.resp
	assert.Equal(t, res.Error,nil, "Expect no Error")

	var missed = 0

	for i := 0; i < MAXITERATIONLIMIT; i++ {
		 if res.Result["Payam"  + strconv.Itoa(i)] != 20 {
			missed++
		 }
	}

	assert.GreaterOrEqual(t, missed, 0,"Can not increment stepcounter before creation")
}

func TestAddGrouprFail(t *testing.T) {
	quit := make(chan bool)
	p := newPedometers("TestPedometers")
	p.startPedometers(quit)
	defer close(quit)
	req := newRequest()
	go p.AddGroup(req)
	res := <- req.resp
	assert.EqualError(t, res.Error,"NO_GROUP","Should generate NO_GROUP")
}


func TestAddGrouprSuccess(t *testing.T) {
	quit := make(chan bool)
	p := newPedometers("TestPedometers")
	p.startPedometers(quit)
	defer close(quit)
	req := newRequest()
	req.Group = "Apsis"
	go p.AddGroup(req)
	res := <- req.resp
	assert.Equal(t, res.Error,nil, "Should generate NO_GROUP")
}


func TestAddGrouprExistFail(t *testing.T) {
	quit := make(chan bool)
	p := newPedometers("TestPedometers")
	p.startPedometers(quit)
	defer close(quit)

	req := newRequest()
	req.Group = "Apsis"
	go p.AddGroup(req)
	_ = <- req.resp

	req = newRequest()
	req.Group = "Apsis"
	go p.AddGroup(req)
	res := <- req.resp

	assert.EqualError(t, res.Error,"GROUP_EXISTS","Should generate GROUP_EXISTS")
}


func TestAddWalkerToGroupSucess(t *testing.T) {

	quit := make(chan bool)
	p := newPedometers("TestPedometers")
	p.startPedometers(quit)
	defer close(quit)

	req := newRequest()
	req.Name = "Payam"
	go p.AddWalker(req)
	res := <- req.resp
	assert.Equal(t, res.Error,nil, "No Error")

	req = newRequest()
	req.Group = "KTH"
	go p.AddGroup(req)
	res = <- req.resp

	assert.Equal(t, res.Error,nil, "No Error")

	req = newRequest()
	req.Name = "Payam"
	req.Group = "KTH"
	go p.AddWalkerToGroup(req)
	res = <- req.resp

	assert.Equal(t, res.Error,nil, "No Error")



	req = newRequest()
	go p.ListAll(req)
	res = <- req.resp
	assert.Equal(t, res.Result["Payam"],0, "Counter set to zer0")
}



func TestAddWalkerToGroupFail(t *testing.T) {

	quit := make(chan bool)
	p := newPedometers("TestPedometers")
	p.startPedometers(quit)
	defer close(quit)

	req := newRequest()
	req.Name = "Payam"
	go p.AddWalker(req)
	res := <- req.resp
	assert.Equal(t, res.Error,nil, "No Error")

	req = newRequest()
	req.Group = "Apsis"
	go p.AddGroup(req)
	_ = <- req.resp

	req = newRequest()
	req.Name = "Payam"
	req.Group = "KTH"
	go p.AddWalkerToGroup(req)
	res = <- req.resp

	assert.EqualError(t, res.Error,"GROUP_MISSING","Should generate GROUP_MISSING")

}



func TestAddWalkerToGroupMultipleGroups(t *testing.T) {

	quit := make(chan bool)
	p := newPedometers("TestPedometers")
	p.startPedometers(quit)
	defer close(quit)

	var waitgroup sync.WaitGroup

	req := newRequest()
	req.Group = "TestGroup"
	go p.AddGroup(req)
	_ = <- req.resp


	for i := 0; i < MAXITERATIONLIMIT; i++ {
		waitgroup.Add(1)
		go func(index int) {
			req := newRequest()
			req.Name = "Payam" + strconv.Itoa(index)
			p.AddWalker(req)
			_ = <-req.resp

			req = newRequest()
			req.Name = "Payam" + strconv.Itoa(index)
			req.Steps = 10
			p.RegisterSteps(req)
			_ = <-req.resp

			req = newRequest()
			req.Name = "Payam" + strconv.Itoa(index)
			req.Group = "TestGroup"
			p.AddWalkerToGroup(req)
			_ = <-req.resp

			waitgroup.Done()
		}(i)
	}


	waitgroup.Wait()

	req = newRequest()
	req.Group = "TestGroup"
	p.ListGroup(req)
	res := <-req.resp

	assert.Equal(t, res.Steps,MAXITERATIONLIMIT * 10, "Testgroup should add up")


	req = newRequest()
	go p.ListAll(req)
	res = <- req.resp
	assert.Equal(t, res.Steps,MAXITERATIONLIMIT * 10, "List all should add up")

}

