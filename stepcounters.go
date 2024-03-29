package main

import (
	"time"
)

const (
	INTERNAL source = iota
	EXTERNAL
)

func (s source) String() string {
	return [...]string{
		"INTERNAL",
		"EXTERNAL",
	}[s]
}

func newPedometers(name string, config *config) *pedometers {
	return &pedometers{
		name,
		make(leaderboard),
		make(groups),
		make(chan *request, config.MAXQUEUELENGTH),
		make(chan *request, config.MAXQUEUELENGTH),
		make(chan *request, config.MAXQUEUELENGTH),
		make(chan *request, config.MAXQUEUELENGTH),
		config,
	}
}

func (p *pedometers) startPedometers(quit chan bool) {
	println("Starting processors")
	go func() {
		println("Starting leaderboard processor")
		for {
			select {
			case req := <-p.leaderBoardCmdInternal:
				println("Processing leaderboard Internal queue", req.Cmd.String())
				p.dispatchCommand(req)
			case req := <-p.leaderBoardCmd:
				println("Processing leaderboard queue", req.Cmd.String())
				p.dispatchCommand(req)
			case <-quit:
				println("Stoping leaderboard processor")
				return
			}
		}
	}()

	go func() {
		println("Starting group processors")
		for {
			select {
			case req := <-p.groupsCmdInternal:
				println("Processing groups internal queue", req.Cmd.String())
				p.dispatchCommand(req)
			case req := <-p.groupsCmd:
				println("Processing groups  queue", req.Cmd.String())
				p.dispatchCommand(req)
			case <-quit:
				println("Stoping group processor")
				return
			}
		}
	}()

}

func (p *pedometers) dispatchCommand(req *request) {
	switch req.Cmd {
	case ADDWALKER:
		p.processAddWalker(req)
	case REGISTERSTEPS:
		p.processRegisterSteps(req)
	case ADDGROUP:
		p.processAddGroup(req)
	case ADDWALKERTOGROUP:
		p.processAddWalkerToGroup(req)
	case DELETEWALKER:
		p.processDeleteWalker(req)
	case GETWALKER:
		p.processGetWalker(req)
	case RESETSTEPS:
		p.processResetSteps(req)
	case LISTGROUP:
		p.processListGroup(req)
	case LISTALL:
		p.processListAll(req)
	case LISTALLGROUPS:
		p.processListAllGroups(req)
	case SCAN:
		p.processScan(req)

	default:
		req.Error = &UnknownCmdError{}
		req.resp <- req
	}
}

func (p *pedometers) execLeadBoardCmd(req *request) {
	waitDuration := time.Duration(p.config.TIMEOUT)
	if (req.Source == EXTERNAL) {
		select {
		case p.leaderBoardCmd <- req:
		case <-time.After(waitDuration * time.Second):
			req.Error = &TimeOutError{}
			req.resp <- req
			close(req.resp)
		}
	} else {
		select {
		case p.leaderBoardCmdInternal <- req:
		case <-time.After(waitDuration  * time.Second):
			req.Error = &TimeOutError{}
			req.resp <- req
			close(req.resp)
		}
	}
}

func (p *pedometers) execGroupCmd(req *request) {
	waitDuration := time.Duration(p.config.TIMEOUT)
	if (req.Source == EXTERNAL) {
		select {
		case p.groupsCmd <- req:
		case <-time.After(waitDuration * time.Second):
			req.Error = &TimeOutError{}
			req.resp <- req
			close(req.resp)
		}
	} else {
		select {
		case p.groupsCmdInternal <- req:
		case <-time.After(waitDuration  * time.Second):
			req.Error = &TimeOutError{}
			req.resp <- req
			close(req.resp)
		}
	}
}
