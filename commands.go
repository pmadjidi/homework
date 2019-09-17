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

func (a *App) steperHash(req *request)  {
	req.Hash = calcHash(req.Name)
}

func (a *App) groupHash(req *request)  {
	req.Hash = calcHash(req.Group)
}

func (a *App) chardName(req *request) string {
	return req.Hash[0:SHARTSLICE]
}





func (a *App) AddWalker(req *request) {
	req.Cmd = ADDWALKER
	if req.Name == EMPTYSTRING {
		req.Error = &InvalidNameError{}
		req.resp <- req
		close(req.resp)
	} else {
		a.execLeadBoardCmd(req)
	}
}

func (a *App) GetWalker(req *request) {
	req.Cmd = GETWALKER
	if req.Name == EMPTYSTRING {
		req.Error = &InvalidNameError{}
		req.resp <- req
		close(req.resp)
	} else {
		a.execLeadBoardCmd(req)
	}
}

func (a *App) RegisterSteps(req *request) {
	req.Cmd = REGISTERSTEPS
	if req.Name == EMPTYSTRING {
		req.Error = &InvalidNameError{}
		req.resp <- req
		close(req.resp)
	} else if req.Steps <= 0 {
		req.Error = &NegativeStepCounterOrZeroError{}
		req.resp <- req
		close(req.resp)
	} else if req.Steps >= a.config.MAXNUMBEROFSTEPSINPUT {
		req.Error = &StepInputOverFlowError{}
		req.resp <- req
		close(req.resp)
	} else {
		a.execLeadBoardCmd(req)
	}
}

func (a *App) AddGroup(req *request) {
	req.Cmd = ADDGROUP
	if req.Group == EMPTYSTRING {
		req.Error = &InvalidGroupNameError{}
		req.resp <- req
		close(req.resp)
	} else {
		a.execGroupCmd(req)
	}
}

func (a *App) AddWalkerToGroup(req *request) {
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
		a.execGroupCmd(req)
	}

}

//not implemented yet
func (a *App) DeleteWalker(req *request) {
	req.Cmd = DELETEWALKER
	req.Error = &NotImplementedError{}
	req.resp <- req
	close(req.resp)
}

func (a *App) ResetSteps(req *request) {
	req.Cmd = RESETSTEPS
	if req.Name == EMPTYSTRING {
		req.Error = &InvalidNameError{}
		req.resp <- req
		close(req.resp)
	} else {
		a.execLeadBoardCmd(req)
	}
}

func (a *App) ListGroup(req *request) {
	req.Cmd = LISTGROUP
	if req.Group == EMPTYSTRING {
		req.Error = &InvalidGroupNameError{}
		req.resp <- req
		close(req.resp)
	} else {
		a.execGroupCmd(req)
	}
}

func (a *App) ListAll(req *request) {
	req.Cmd = LISTALL
	a.execLeadBoardCmd(req)
}

func (a *App) ListAllGroups(req *request) {
	req.Cmd = LISTALLGROUPS
	a.execLeadBoardCmd(req)
}

func (a *App) scan(req *request) {
	req.Cmd = SCAN
	a.execLeadBoardCmd(req)
}
