package main

import "sync"
import "fmt"



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
			println("name is",name)
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
	println("finishing processListGroups:",p.index)
}

func (p *pedometers) processListShardGroups(req *request) {
	println("processListAllGroupsNode begin")
	req.Results = make(map[string]leaderboard)
	for k, _ := range p.groups {
		newRequest := newRequestInternal()
		newRequest.Group = k
		APP.GetGroup(newRequest)
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
	println("processListAllGroupsNode end")
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
				newreq.shard = APP.shards[index]
				APP.ListGroups(newreq)
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
		select {
		case r := <-resp:

			fmt.Printf("%+v\n", r)
			println("&&&&&&& r is", r, r.Hash)
			if r.Error == nil {
				for k, v := range r.Results {
					req.Results[k] = v
					println("%%%%%", k, v)
				}
			} else {
				println("Recived: Error", r.Error.Error())
				req.Error = r.Error
				req.resp <- req
				close(req.resp)
				return
			}
		default:

		}
	}



	/*
		println("Here")
		localReq := newRequestInternal()
		p.processListGroups(localReq)
		resp := <-localReq.resp
		println("Here after processListGroups")

		for k, v := range resp.Results {
			req.Results[k] = v
		}

	*/
	println("almost")

	req.resp <- req
	println("at")
	close(req.resp)
	println("end")
}
