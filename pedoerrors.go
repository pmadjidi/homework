package main

import (
	"fmt"
)

type TimeOutError struct {
}

func (e *TimeOutError) Error() string {
	return fmt.Sprintf("TIMEOUT")
}


type InvalidNameError struct {
}

func (e *InvalidNameError) Error() string {
	return fmt.Sprintf("NO_NAME")
}

type NameExistsError struct {
}

func (e *NameExistsError) Error() string {
	return fmt.Sprintf("NAME_EXISTS")
}

type NameDoesNotExistsError struct {
}

func (e *NameDoesNotExistsError) Error() string {
	return fmt.Sprintf("NAME_MISSING")
}


type InvalidGroupNameError struct {
}

func (e *InvalidGroupNameError) Error() string {
	return fmt.Sprintf("NO_GROUP")
}



type GroupExistsError struct {
}

func (e *GroupExistsError) Error() string {
	return fmt.Sprintf("GROUP_EXISTS")
}

type GroupDoesNotExistsError struct {
}

func (e *GroupDoesNotExistsError) Error() string {
	return fmt.Sprintf("GROUP_MISSING")
}

type NotImplementedError struct {
}

func (e *NotImplementedError) Error() string {
	return fmt.Sprintf("COMMAND_NOT_IMPLEMENTED")
}

type UnknownCmdError struct {
}

func (e *UnknownCmdError) Error() string {
	return fmt.Sprintf("UNKOWN_CMD")
}
