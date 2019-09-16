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

type NegativeStepCounterOrZeroError struct {
}

func (e *NegativeStepCounterOrZeroError) Error() string {
	return fmt.Sprintf("NEGATIVE_STEP_COUNTER_OR_ZERO")
}

type StepInputOverFlowError struct {
}

func (e *StepInputOverFlowError) Error() string {
	return fmt.Sprintf("STEP_INPUT_OVERFLOW")
}

type MaxNumberOFWalkersReachedError struct {
}

func (e *MaxNumberOFWalkersReachedError) Error() string {
	return fmt.Sprintf("MAX_NUMBER_OF_WALKERS_REACHED")
}

type MaxNumberOFGroupsReachedError struct {
}

func (e *MaxNumberOFGroupsReachedError) Error() string {
	return fmt.Sprintf("MAX_NUMBER_OF_GROUPS_REACHED")
}

type MaxNumberOFWalkersInGroupsReachedError struct {
}

func (e *MaxNumberOFWalkersInGroupsReachedError) Error() string {
	return fmt.Sprintf("MAX_NUMBER_OF_WALKERS_IN_GROUP_REACHED")
}

type EnvVariableMissingOrNotIntError struct {
}

func (e *EnvVariableMissingOrNotIntError) Error() string {
	return fmt.Sprintf("ENV-VARIABLE-MISSING-OR-NOT-INT")
}
