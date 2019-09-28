package main



const (
	NOP command = iota
	ADDUSER
	GETUSER
	DELETEUSER
	USERSFORSHARD
	USERS
	REGISTERPOINTS
	RESETPOINTS
	ADDGROUP
	ADDUSERTOGROUP
	GETGROUP
	GROUPS
	GROUPSFORASHARD
)

func (c command) String() string {
	return [...]string{
		"NOP",
		"ADDUSER",
		"GETUSER",
		"DELETEUSER",
		"USERSFORSHARD",
		"USERS",
		"REGISTERPOINTS",
		"RESETPOINTS",
		"ADDGROUP",
		"ADDUSERTOGROUP",
		"GETGROUP",
		"GROUPS",
		"GROUPSFORASHARD",
	}[c]
}

func (a *App) userHash(req *request)  {
	cachhit,found := a.cache.Load(req.Name)
	if !found {
	req.Hash = calcHash(req.Name)
	a.cache.Store(req.Name ,req.Hash)
	} else {
		req.Hash = cachhit.(string)
	}
}

func (a *App) groupHash(req *request)  {
	cachhit,found := a.cache.Load(req.Group)
	if !found {
		req.Hash = calcHash(req.Group)
		a.cache.Store(req.Group ,req.Hash)
	} else {
		req.Hash = cachhit.(string)
	}
}

func (a *App) chardName(req *request) string {
	return req.Hash[0:a.config.SHARDS]
}





func (a *App) AddUser(req *request) {
	req.Cmd = ADDUSER
	if req.Name == EMPTYSTRING {
		req.Error = &InvalidNameError{}
		req.resp <- req
		close(req.resp)
	} else {
		a.execLeadBoardCmd(req)
	}
}

func (a *App) GetUser(req *request) {
	req.Cmd = GETUSER
	if req.Name == EMPTYSTRING {
		req.Error = &InvalidNameError{}
		req.resp <- req
		close(req.resp)
	} else {
		a.execLeadBoardCmd(req)
	}
}

func (a *App) RegisterPoints(req *request) {
	req.Cmd = REGISTERPOINTS
	if req.Name == EMPTYSTRING {
		req.Error = &InvalidNameError{}
		req.resp <- req
		close(req.resp)
	} else if req.Points <= 0 {
		req.Error = &NegativeStepCounterOrZeroError{}
		req.resp <- req
		close(req.resp)
	} else if req.Points >= a.config.MAXNUMBEROFSTEPSINPUT {
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

func (a *App) AddUserToGroup(req *request) {
	req.Cmd = ADDUSERTOGROUP
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
func (a *App) DeleteUser(req *request) {
	req.Cmd = DELETEUSER
	req.Error = &NotImplementedError{}
	req.resp <- req
	close(req.resp)
}

func (a *App) ResetSteps(req *request) {
	req.Cmd = RESETPOINTS
	if req.Name == EMPTYSTRING {
		req.Error = &InvalidNameError{}
		req.resp <- req
		close(req.resp)
	} else {
		a.execLeadBoardCmd(req)
	}
}

func (a *App) GetGroup(req *request) {
	req.Cmd = GETGROUP
	if req.Group == EMPTYSTRING {
		req.Error = &InvalidGroupNameError{}
		req.resp <- req
		close(req.resp)
	} else {
		a.execGroupCmd(req)
	}
}

func (a *App) UsersForShard(req *request) {
	req.Cmd = USERSFORSHARD
	a.execLeadBoardCmd(req)
}

func (a *App) GetUsers(req *request) {
	req.Cmd = GETUSER
	req.Name = RandomString(10)
	a.execLeadBoardCmd(req)
}



func (a *App) GroupsForAShard(req *request) {
	req.Cmd = GROUPSFORASHARD
	a.execGroupCmd(req)
}



func (a *App) Groups(req *request) {
	req.Cmd = GROUPS
	req.Group = RandomString(10)
	a.execGroupCmd(req)
}

