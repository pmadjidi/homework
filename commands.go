package main

const (
	NOP command = iota
	ADDWALKER
	REGISTERSTEPS
	GETWALKER
	ADDGROUP
	ADDWALKERTOGROUP
	DELETEWALKER
	RESETSTEPS
	LISTGROUP
	LISTALL
	LISTALLGROUPS
	SCAN
)

func (c command) String() string {
	return [...]string{
		"NOP",
		"ADDWALKER",
		"REGISTERSTEPS",
		"GETWALKER",
		"ADDGROUP",
		"ADDWALKERTOGROUP",
		"DELETEWALKER",
		"RESETSTEPS",
		"LISTGROUP",
		"LISTALL",
		"LISTALLGROUPS",
		"SCAN",
	}[c]
}

func (p *pedometers) AddWalker(req *request) {
	req.Cmd = ADDWALKER
	if req.Name == EMPTYSTRING {
		req.Error = &InvalidNameError{}
		req.resp <- req
		close(req.resp)
	} else {
		p.execLeadBoardCmd(req)
	}
}

func (p *pedometers) GetWalker(req *request) {
	req.Cmd = GETWALKER
	if req.Name == EMPTYSTRING {
		req.Error = &InvalidNameError{}
		req.resp <- req
		close(req.resp)
	} else {
		p.execLeadBoardCmd(req)
	}
}

func (p *pedometers) RegisterSteps(req *request) {
	req.Cmd = REGISTERSTEPS
	if req.Name == EMPTYSTRING {
		req.Error = &InvalidNameError{}
		req.resp <- req
		close(req.resp)
	} else if req.Steps <= 0 {
		req.Error = &NegativeStepCounterOrZeroError{}
		req.resp <- req
		close(req.resp)
	} else if req.Steps >= p.config.MAXNUMBEROFSTEPSINPUT {
		req.Error = &StepInputOverFlowError{}
		req.resp <- req
		close(req.resp)
	} else {
		p.execLeadBoardCmd(req)
	}
}

func (p *pedometers) AddGroup(req *request) {
	req.Cmd = ADDGROUP
	if req.Group == EMPTYSTRING {
		req.Error = &InvalidGroupNameError{}
		req.resp <- req
		close(req.resp)
	} else {
		p.execGroupCmd(req)
	}
}

func (p *pedometers) AddWalkerToGroup(req *request) {
	req.Cmd = ADDWALKERTOGROUP
	if req.Name == EMPTYSTRING {
		req.Error = &InvalidNameError{}
		req.resp <- req
		close(req.resp)
	} else if req.Group == EMPTYSTRING {
		req.Error = &InvalidGroupNameError{}
		req.resp <- req
		close(req.resp)
	} else {
		p.execGroupCmd(req)
	}

}

//not implemented yet
func (p *pedometers) DeleteWalker(req *request) {
	req.Cmd = DELETEWALKER
	req.Error = &NotImplementedError{}
	req.resp <- req
	close(req.resp)
}

func (p *pedometers) ResetSteps(req *request) {
	req.Cmd = RESETSTEPS
	if req.Name == EMPTYSTRING {
		req.Error = &InvalidNameError{}
		req.resp <- req
		close(req.resp)
	} else {
		p.execLeadBoardCmd(req)
	}
}

func (p *pedometers) ListGroup(req *request) {
	req.Cmd = LISTGROUP
	if req.Group == EMPTYSTRING {
		req.Error = &InvalidGroupNameError{}
		req.resp <- req
		close(req.resp)
	} else {
		p.execGroupCmd(req)
	}
}

func (p *pedometers) ListAll(req *request) {
	req.Cmd = LISTALL
	p.execLeadBoardCmd(req)
}

func (p *pedometers) scan(req *request) {
	req.Cmd = SCAN
	p.execLeadBoardCmd(req)
}
