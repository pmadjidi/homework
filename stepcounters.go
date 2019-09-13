package main

import (
	"time"
)

func newPedometers(name string) *pedometers {
	return &pedometers{
		name,
		make(leaderboard),
		make(groups),
		make(chan *request, MAXQUEUELENGTH),
		make(chan *request, MAXQUEUELENGTH),
	}
}

func (p *pedometers) startPedometers(quit chan bool) {
	println("Starting processors")
	go func() {
		println("Starting leaderboard processor")
		for {
			select {
			case req := <-p.leaderBoardCmd:
				println("Processing", req.Cmd.String())
				p.dispatchCommand(req)
			case <-quit:
				println("Stoping leaderboard processor")
				return
			default:
			}
		}
	}()

	go func() {
		println("Starting group processors")
		for {
			select {
			case req := <-p.groupsCmd:
				println("Processing", req.Cmd.String())
				p.dispatchCommand(req)
			case <-quit:
				println("Stoping group processor")
				return
			default:
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
	select {
	case p.leaderBoardCmd <- req:
	case <-time.After(TIMEOUT * time.Second):
		req.Error = &TimeOutError{}
		req.resp <- req
	}
}

func (p *pedometers) execGroupCmd(req *request) {
	select {
	case p.groupsCmd <- req:
	case <-time.After(TIMEOUT * time.Second):
		req.Error = &TimeOutError{}
		req.resp <- req
	}
}
