package main

import "strconv"

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

func newPedometers(name int, config *config) *pedometers {
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
	name := strconv.Itoa(p.index)
	println("Starting processors",name)
	go func() {
		println("Starting leaderboard processor",name)
		for {
			select {
			case req := <-p.leaderBoardCmdInternal:
				println("Processing leaderboard Internal queue",name, req.Cmd.String())
				p.dispatchCommand(req)
			case req := <-p.leaderBoardCmd:
				println("Processing leaderboard queue", name,req.Cmd.String())
				p.dispatchCommand(req)
			case <-quit:
				println("Stoping leaderboard processor",name)
				return
			}
		}
	}()

	go func() {
		println("Starting group processors",name)
		for {
			select {
			case req := <-p.groupsCmdInternal:
				println("Processing groups internal queue",name, req.Cmd.String())
				p.dispatchCommand(req)
			case req := <-p.groupsCmd:
				println("Processing groups  queue",name, req.Cmd.String())
				p.dispatchCommand(req)
			case <-quit:
				println("Stoping group processor",name)
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
	case GETGROUP:
		p.processGetGroup(req)
	case LISTWALKERS:
		p.processListWalkers(req)
	case LISTALLWALKERS:
		p.processListAllWalkers(req)
	case LISTALLGROUPS:
		p.processListAllGroups(req)
	case LISTGROUPS:
		p.processListGroups(req)
	case LISTGROUPSFORASHARD:
		p.processListShardGroups(req)
	case LISTGROUPSFORASHARDSEQ:
		p.processListShardGroupsSeq(req)
	default:
		req.Error = &UnknownCmdError{}
		req.resp <- req
	}
}

