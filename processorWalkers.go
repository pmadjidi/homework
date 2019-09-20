package main

import (
	"sync"
)

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


func (p *pedometers) processAddWalkerToGroup(req *request) {
	_, groupfound := p.groups[req.Group]

	if !groupfound {
		req.Error = &GroupDoesNotExistsError{}
	} else {

		newReq := newRequestInternal()
		newReq.Name = req.Name
		APP.GetWalker(newReq)
		newResp := <-newReq.resp
		if newResp.Error != nil {
			req.Error = newResp.Error
		} else {
			p.groups[req.Group][req.Name] = 1
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
	close(req.resp)
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

func (p *pedometers) processListAllWalkers(req *request) {

	req.Result = make(leaderboard)
	var wg sync.WaitGroup
	var responseFromOthers = make(chan chan *request, p.config.SHARDS -1)
	for shard := 0; shard < p.config.SHARDS ; shard++ {
		if shard != p.index { // Obs Important to avoid dealock....
			wg.Add(1)
			go func(index int) {
				println("shard", index)
				newreq := newRequestInternal()
				newreq.index = index
				APP.ListWalkers(newreq)
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
		if r.Error == nil {
			for k, v := range r.Result {
				req.Result[k] = v
			}
			req.Steps += r.Steps

		} else {
			println("Recived: Error", r.Error.Error())
			req.Error = r.Error
			req.resp <- req
			close(req.resp)
			return
		}
	}

	for k, v := range p.leaderboard {
		req.Result[k] = v
	}

	req.resp <- req
	close(req.resp)
}

func (p *pedometers) processListWalkers(req *request) {
	shard := req.index
	println("processListWalkers", shard)
	req.Result = make(leaderboard)
	req.Steps = 0
	for k, v := range p.leaderboard {
		req.Result[k] = v
		req.Steps += v
	}

	req.resp <- req
	close(req.resp)
}

