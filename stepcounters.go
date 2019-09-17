package main

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
	println("Starting processors",p.name)
	go func() {
		println("Starting leaderboard processor",p.name)
		for {
			select {
			case req := <-p.leaderBoardCmdInternal:
				println("Processing leaderboard Internal queue",p.name, req.Cmd.String())
				p.dispatchCommand(req)
			case req := <-p.leaderBoardCmd:
				println("Processing leaderboard queue", p.name,req.Cmd.String())
				p.dispatchCommand(req)
			case <-quit:
				println("Stoping leaderboard processor",p.name)
				return
			}
		}
	}()

	go func() {
		println("Starting group processors",p.name)
		for {
			select {
			case req := <-p.groupsCmdInternal:
				println("Processing groups internal queue",p.name, req.Cmd.String())
				p.dispatchCommand(req)
			case req := <-p.groupsCmd:
				println("Processing groups  queue",p.name, req.Cmd.String())
				p.dispatchCommand(req)
			case <-quit:
				println("Stoping group processor",p.name)
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

