package main

import "sync"

func (p *pedometers) processAddGroup(req *request) {
	_, found := p.groups[req.Group]
	if found {
		req.Error = &GroupExistsError{}
	} else if len(p.groups) >= p.config.MAXNUMBEROFGROUPS {
		req.Error = &MaxNumberOFGroupsReachedError{}
	} else {
		p.groups[req.Group] = make(leaderboard)
	}
	req.resp <- req
	close(req.resp)
}

func (p *pedometers) processGetGroup(req *request) {
	req.Result = make(leaderboard)
	aGroup, found := p.groups[req.Group]
	if !found {
		req.Error = &GroupDoesNotExistsError{}
		req.resp <- req
		close(req.resp)
		return
	} else {
		var wg sync.WaitGroup
		var responseFromOthers = make(chan chan *request, len(aGroup))
		for k, _ := range aGroup {
			req.Result[k] = 0
			wg.Add(1)
			go func(name string) {
				newreq := newRequestInternal()
				newreq.Name = name
				println("Name is:", newreq.Name)
				APP.GetWalker(newreq)
				responseFromOthers <- newreq.resp
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

func (p *pedometers) processListGroups(req *request) {
	println("begin1")
	req.Results = make(map[string]leaderboard)
	var wg sync.WaitGroup
	var responseFromOthers = make(chan chan *request, len(p.groups))
	for k, _ := range p.groups {
		wg.Add(1)
		go func(name string) {
			println("name is", name)
			newreq := newRequestInternal()
			newreq.Group = name
			APP.GetGroup(newreq)
			responseFromOthers <- newreq.resp
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
			r.Result["TOTAL: "+r.Group ] = r.Steps
			req.Results[r.Group] = r.Result
			PrettyPrint(req.Results)
		} else {
			req.Error = r.Error
			req.resp <- req
			close(req.resp)
			return
		}
	}

	req.resp <- req
	close(req.resp)
	println("finishing processListGroups:", p.index)
}

func (p *pedometers) processListShardGroups(req *request) {
	println("processListAllGroupsNode begin: ", req.index)
	req.Results = make(map[string]leaderboard)
	var wg sync.WaitGroup
	responseFromOthers := make(chan chan *request, len(p.groups))
	for aGroupName, aGroup := range p.groups {
		groupReplica := make(leaderboard)
		for aPerson, _ := range aGroup {
			wg.Add(1)
			go func(name, group string) {
				newRequest := newRequestInternal()
				newRequest.Name = name
				newRequest.Group = group
				APP.GetWalker(newRequest)
				responseFromOthers <- newRequest.resp
				wg.Done()
			}(aPerson, aGroupName)
		}
		req.Results[aGroupName] = groupReplica
	}

	go func() {
		wg.Wait()
		close(responseFromOthers)
	}()

	for resp := range responseFromOthers {
		r := <-resp
		if r.Error == nil {
			req.Results[r.Group][r.Name] = r.Steps
		} else {
			println(r.Error.Error())
		}
	}

	req.resp <- req
	close(req.resp)
	println("processListAllGroupsNode end: ", req.index)
}

func (p *pedometers) processListShardGroupsSeq(req *request) {
	println("processListAllGroupsNode begin: ", req.index)
	req.Results = make(map[string]leaderboard)
	for aGroupName, aGroup := range p.groups {
		groupReplica := make(leaderboard)
		groupTotal := 0
		for aPerson, _ := range aGroup {
			newRequest := newRequestInternal()
			newRequest.Name = aPerson
			APP.GetWalker(newRequest)
			newResp := <-newRequest.resp
			if newResp.Error == nil {
				groupReplica [aPerson] = newResp.Steps
				groupTotal +=  newResp.Steps
			} else {
				req.Error = newResp.Error
				req.resp <- req
				close(req.resp)
				return
			}
			groupReplica["-SUM:"] = groupTotal
			req.Results[aGroupName] = groupReplica

		}
	}

	req.resp <- req
	close(req.resp)
	println("processListAllGroupsNode end: ", req.index)
}

func (p *pedometers) processListAllGroups(req *request) {
	req.Results = make(map[string]leaderboard)
	var wg sync.WaitGroup

	var responseFromOthers = make(chan chan *request, p.config.NUMBEROFSHARDS)
	for shard := 0; shard < p.config.NUMBEROFSHARDS; shard++ {
		if shard != p.index { // Obs Important to avoid deadlock....
			wg.Add(1)
			go func(index int) {
				newreq := newRequestInternal()
				newreq.index = index
				APP.ListGroupsForAShardSeq(newreq)
				responseFromOthers <- newreq.resp
				wg.Done()
			}(shard)
		}
	}

	go func() {
		wg.Wait()
		close(responseFromOthers)
	}()

	for resp := range responseFromOthers {
		r := <-resp
		for k, v := range r.Results {
			req.Results[k] = v
		}
	}

	for aGroupName, aGroup := range p.groups {
		groupReplica := make(leaderboard)
		groupTotal := 0
		for aPerson, _ := range aGroup {
			newRequest := newRequestInternal()
			newRequest.Name = aPerson
			APP.GetWalker(newRequest)
			newResp := <-newRequest.resp
			if newResp.Error == nil {
				groupReplica [aPerson] = newResp.Steps
				groupTotal +=  newResp.Steps
			} else {
				req.Error = newResp.Error
				req.resp <- req
				close(req.resp)
				return
			}
			groupReplica["-SUM"] = groupTotal
			req.Results[aGroupName] = groupReplica
		}
	}
	req.resp <- req
	close(req.resp)
}
