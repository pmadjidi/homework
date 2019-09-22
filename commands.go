package main

import (
	"sync"
	"sync/atomic"
)


func (p *pedometers) AddWalker(name string) error {
	var err error = nil
	var steps int32
	if name == EMPTYSTRING {
		return &InvalidNameError{}
	} else {
		_, ok := p.leaderboard.Load(name)
		if ok {
			err = &NameExistsError{}
		} else {
			steps = 0
			p.leaderboard.Store(name, &steps)
		}
		return err
	}
}

func (p *pedometers) GetWalker(name string) (steps int, err error) {
	steps = -1
	err = nil
	var stepPointer *int32 = nil
	if name == EMPTYSTRING {
		err = &InvalidNameError{}
	} else {
		val, ok := p.leaderboard.Load(name)
		if !ok {
			err = &NameDoesNotExistsError{}
		} else {
			stepPointer = val.(*int32)
			steps = int(*stepPointer)
		}
	}
	return
}

func (p *pedometers) RegisterSteps(name string, incSteps int32) (steps int32, err error) {
	steps = -1
	err = nil
	var stepPointer *int32 = nil
	if name == EMPTYSTRING {
		err = &InvalidNameError{}
	} else if incSteps <= 0 {
		err = &NegativeStepCounterOrZeroError{}
	} else if incSteps >= int32(p.config.MAXNUMBEROFSTEPSINPUT) {
		err = &StepInputOverFlowError{}
	} else {
		val, ok := p.leaderboard.Load(name)
		if !ok {
			err = &NameDoesNotExistsError{}
		} else {
			stepPointer = val.(*int32)
			steps = atomic.AddInt32(stepPointer, incSteps)
		}
	}
	return
}


func (p *pedometers) AddGroup(name string) error {
	var err error = nil
	if name == EMPTYSTRING {
		return &InvalidNameError{}
	} else {
		_, ok := p.groups.Load(name)
		if ok {
			err = &GroupExistsError{}
		} else {
			var aGroup sync.Map
			p.groups.Store(name, &aGroup)
		}
		return err
	}
}

func (p *pedometers) getGroup(name string) (*sync.Map,error) {
	var err error = nil
	var group *sync.Map = nil
	if name == EMPTYSTRING {
		err = &InvalidNameError{}
	} else {
		aGroup , ok := p.groups.Load(name)
		if !ok {
			err = &GroupDoesNotExistsError{}
		} else {
			group = aGroup.(*sync.Map)
		}
	}
	return group,err
}

func (p *pedometers) AddWalkerToGroup(name,group string) error {
	var ERR error = nil
	if name  == EMPTYSTRING {
		ERR = &InvalidNameError{}
	} else if group == EMPTYSTRING {
		ERR = &InvalidGroupNameError{}
	} else {
		if _, err := p.GetWalker(name); err != nil && err.Error() == "NAME_MISSING" {
			ERR = &NameDoesNotExistsError{}
		} else if agroup, err := p.getGroup(group); err == nil {
			agroup.Store(name, true)
		} else {
			ERR =  err
		}
	}
	return ERR
}

//not implemented yet
func (p *pedometers) DeleteWalker(name string)error {
	return &NotImplementedError{}

}

func (p *pedometers) ResetSteps(name string) error {
	return &NotImplementedError{}
}

func (p *pedometers) ListGroup(group string) (map[string]int,error) {
	var ERR error = nil
	var foundGroup = make(map[string]int)
	foundGroup["TOTAL"] = 0
	if group  == EMPTYSTRING {
		ERR  = &InvalidGroupNameError{}
	} else {
		aGroup,err := p.getGroup(group)
		if err != nil {
			ERR = err
			println("ListGroup Error")
		} else {
			aGroup.Range(func(k,v interface{}) bool {
				key := k.(string)
				steps, e := p.GetWalker(key)
				if e != nil {
					ERR = e
					return false
				} else {
					foundGroup[key] = steps
					foundGroup["TOTAL"] += steps
				}
				return true
			})
		}
	}
	return foundGroup,ERR
}

func (p *pedometers) ListAllSteppers() map[string]int {
	var AllSteppers = make(map[string]int)
	var total int = 0
	p.leaderboard.Range(func(k,v interface{}) bool {
		key := k.(string)
		val := *v.(*int32)
		AllSteppers[key] = int(val)
		total += int(val)
		return true
	})
	AllSteppers["TOTAL"] = total
	return AllSteppers
}

func (p *pedometers) ListAllGroups() (map[string]map[string]int,error) {
	var err error = nil
	AllGroups := make(map[string]map[string]int)
	p.groups.Range(func(k,v interface{}) bool {
		aGroup,e  := p.ListGroup(k.(string))
		if e != nil {
			err = e
			return false
		}
		AllGroups[k.(string)] = aGroup
		return true
	})

	return AllGroups,err

}
