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
				APP.GetUser(newreq)
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
				req.Result[r.Name] = r.Points
				req.Points += r.Points
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

func (p *pedometers) processGroupForShards(req *request) {
	req.Results = make(map[string]leaderboard)
	for aGroupName, aGroup := range p.groups {
		groupReplica := make(leaderboard)
		req.Results[aGroupName] = groupReplica
		groupReplica["SUM"] = 0
		for aPerson, _ := range aGroup {
			newRequest := newRequestInternal()
			newRequest.Name = aPerson
			APP.GetUser(newRequest)
			newResp := <-newRequest.resp
			if newResp.Error == nil {
				groupReplica [aPerson] = newResp.Points
				groupReplica["SUM"] += newResp.Points
			} else {
				req.Error = newResp.Error
				req.resp <- req
				close(req.resp)
				return
			}
		}
	}

	if (req.Results != nil ) {
		PrettyPrint(req.Results)
	}
	req.resp <- req
	close(req.resp)
}

func (p *pedometers) processGroups(req *request) {
	req.Results = make(map[string]leaderboard)

	for shard := 0; shard < p.config.SHARDS; shard++ {
		if shard != p.index { // Obs Important to avoid deadlock....
			newreq := newRequestInternal()
			newreq.index = shard
			APP.GroupsForAShard(newreq)
			responseFromOthersShards := <-newreq.resp
			if responseFromOthersShards.Error != nil {
				req.resp <- req
				close(req.resp)
				return
			} else {
				if responseFromOthersShards.Results != nil {
					for k, v := range responseFromOthersShards.Results {
						req.Results[k] = v
					}
				}
			}
		}
	}

	for aGroupName, aGroup := range p.groups {
		groupReplica := make(leaderboard)
		groupTotal := 0
		for aPerson, _ := range aGroup {
			newRequest := newRequestInternal()
			newRequest.Name = aPerson
			APP.GetUser(newRequest)
			newResp := <-newRequest.resp
			if newResp.Error == nil {
				groupReplica [aPerson] = newResp.Points
				groupTotal += newResp.Points
			} else {
				req.Error = newResp.Error
				req.resp <- req
				close(req.resp)
				return
			}
			groupReplica["SUM"] = groupTotal
			req.Results[aGroupName] = groupReplica

		}
	}

	req.resp <- req
	close(req.resp)
}
