package main

import "sync"

func (p *pedometers) processAddWalker(req *request) {

	_, found := p.leaderboard[req.Name]
	if found {
		req.Error = &NameExistsError{}
	} else if len(p.leaderboard) >= p.config.MAXNUMBERSOFWALKERS {
		req.Error = &MaxNumberOFWalkersReachedError{}
	} else {
		p.leaderboard[req.Name] = 0
	}
	req.resp <- req
	close(req.resp)
}

func (p *pedometers) processGetWalker(req *request) {
	steps, found := p.leaderboard[req.Name]
	if !found {
		req.Error = &NameDoesNotExistsError{}
	} else {
		req.Steps = steps
	}
	req.resp <- req
	close(req.resp)
}

func (p *pedometers) processRegisterSteps(req *request) {
	_, found := p.leaderboard[req.Name]
	if !found {
		req.Error = &NameDoesNotExistsError{}
	} else {
		p.leaderboard[req.Name] += req.Steps
		req.Steps = p.leaderboard[req.Name]
	}
	req.resp <- req
	close(req.resp)
}

func (p *pedometers) processAddGroup(req *request) {
	_, found := p.groups[req.Group]
	if found {
		req.Error = &GroupExistsError{}
	} else if len(p.groups) >= p.config.MAXNUMBEROFGROUPS {
		req.Error = &MaxNumberOFGroupsReachedError{}
	} else {
		p.groups[req.Group] = make(map[string]bool)
	}
	req.resp <- req
	close(req.resp)
}

func (p *pedometers) processAddWalkerToGroup(req *request) {
	aUserGroup, groupfound := p.groups[req.Group]

	if !groupfound {
		req.Error = &GroupDoesNotExistsError{}
	} else if _, userfound := p.groups[req.Group][req.Name]; userfound {
		req.Error = &NameExistsError{}
	} else if len(aUserGroup) >= p.config.MAXNUMBEROFWALKERSINGROUP {
		req.Error = &MaxNumberOFWalkersInGroupsReachedError{}
	} else {

		newReq := newRequestInternal()
		newReq.Name = req.Name
		APP.GetWalker(newReq)
		newResp := <-newReq.resp
		if newResp.Error != nil {
			req.Error = newResp.Error
		} else {
			p.groups[req.Group][req.Name] = true
			req.Steps = newResp.Steps
		}
	}
	req.resp <- req
	close(req.resp)
}

//not implemented yet...
func (p *pedometers) processDeleteWalker(req *request) {
	req.Error = &NotImplementedError{}
	req.resp <- req
}

func (p *pedometers) processResetSteps(req *request) {
	_, found := p.leaderboard[req.Name]
	if !found {
		req.Error = &NameDoesNotExistsError{}
	} else {
		p.leaderboard[req.Name] = 0
	}
	req.resp <- req
	close(req.resp)
}



func (p *pedometers) processListGroup(req *request) {
	aGroup, found := p.groups[req.Group]
	if !found {
		req.Error = &GroupDoesNotExistsError{}
	} else {

		req.Result = make(leaderboard)
		var wg sync.WaitGroup
		var responseFromOthers = make(chan chan *request, len(aGroup))
		for k, _ := range aGroup {
			req.Result[k] = 0
			wg.Add(1)
			go func(name string) {
				req := newRequestInternal()
				req.Name = k
				APP.GetWalker(req)
				responseFromOthers <- req.resp
				wg.Done()
			}(k)
		}

		go func() {
			wg.Wait()
			close(responseFromOthers)
		}()

		for resp := range responseFromOthers {
			r := <-resp
			if r.Error == nil {
				req.Result[r.Name] = r.Steps
				req.Steps += r.Steps
			} else {
				req.Error = r.Error
				req.resp <- req
				close(req.resp)
				return
			}
		}

		req.resp <- req
		close(req.resp)
	}
}

func (p *pedometers) processListAll(req *request) {
	result := make(leaderboard)
	req.Steps = 0
	for k, v := range p.leaderboard {
		result[k] = p.leaderboard[k]
		req.Steps += v
	}
	req.Result = result
	req.resp <- req
	close(req.resp)
}

func (p *pedometers) processScan(req *request) {
	if req.Group == EMPTYSTRING {
		for k, v := range p.leaderboard {
			req.Result[k] = v
		}
	} else {
		for k, _ := range req.Result {
			steps, found := p.leaderboard[k]
			if found {
				req.Result[k] = steps
			} else {
				req.Result[k] = NOTFOUND // prepare for implementation of delte function
			}
		}
	}
	req.resp <- req
	close(req.resp)
}

func (p *pedometers) processListAllGroups(req *request) {
	req.Results = make(map[string]leaderboard)
	for k, _ := range p.groups {
		newRequest := newRequestInternal()
		newRequest.Group = k
		APP.ListGroup(newRequest)
		ans := <-newRequest.resp
		if ans.Error != nil {
			req.Error = newRequest.Error
			req.resp <- req
			return
		} else {

			ans.Result[k+"-total"] = ans.Steps
			req.Results[k] = ans.Result
		}
	}
	req.resp <- req
	close(req.resp)

}
