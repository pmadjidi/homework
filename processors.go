package main



func (p *pedometers) processAddWalker(req *request) {
	_, found := p.leaderboard[req.Name]
	if found {
		req.Error = &NameExistsError{}
	} else {
		p.leaderboard[req.Name] = 0
	}
	req.resp <- req
}



func (p *pedometers) processGetWalker(req *request) {
	steps, found := p.leaderboard[req.Name]
	if !found {
		req.Error = &NameDoesNotExistsError{}
	} else {
		req.Steps = steps
	}
	req.resp <- req
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
}



func (p *pedometers) processAddGroup(req *request) {
	_, found := p.groups[req.Group]
	if found {
		req.Error = &GroupExistsError{}
	} else {
		p.groups[req.Group] = make(map[string]bool)
	}
	req.resp <- req
}





func (p *pedometers) processAddWalkerToGroup(req *request) {
	_, groupfound := p.groups[req.Group]

	if !groupfound {
		req.Error = &GroupDoesNotExistsError{}
		req.resp <- req
		return
	}

	_, userfound := p.groups[req.Group][req.Name]

	if userfound {
		req.Error = &NameExistsError{}
		req.resp <- req
		return
	} else {
		newReq := newRequest()
		newReq.Name = req.Name
		p.GetWalker(newReq)
		newResp := <- newReq.resp
		if newResp.Error != nil {
			req.Error = newResp.Error
		} else {
			p.groups[req.Group][req.Name] = true
			req.Steps = newResp.Steps
		}
	}
	req.resp <- req
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
}



func (p *pedometers) processListGroup(req *request) {
	aGroup, found := p.groups[req.Group]
	if !found {
		req.Error = &GroupDoesNotExistsError{}
	} else {
		newReq := newRequest()
		newReq.Group = req.Group
		 newReq.Result = make(leaderboard)
		for k, _ := range aGroup {
			newReq.Result[k] = 0
		}
		p.scan(newReq)
		newResp := <- newReq.resp
		if newResp.Error != nil {
			req.Error = newResp.Error
		} else {
			req.Result = newReq.Result
			for k,v := range req.Result {
				if v == NOTFOUND { // prepare for implementation of delte function
					delete(req.Result,k)
				} else {
					req.Steps += v
				}
			}
		}
	}
	req.resp <- req
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
}

func (p *pedometers) processListAllGroups(req *request) {

	for k, _ := range p.groups {
		newRequest := newRequest()
		newRequest.Group = k
		p.ListGroup(newRequest)
		ans := <-newRequest.resp
		if ans.Error != nil {
			req.Error = newRequest.Error
			req.resp <- req
			return
		} else {

			ans.Result[k] = ans.Steps
			req.Results[k] = ans.Result
		}
	}
	req.resp <- req

}

