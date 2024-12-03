package base

import (
	"fmt"
)

type AtomTask interface {
	Name() string
	Id() int8
	Execute(CtxStorage) error
	Rollback(CtxStorage) error
	SetErrors(isExecute bool, f, t int8, fName, tName string, err error)
	GetErrors() (*ErrExecute, *ErrRollback)
}

type TaskFunc func(CtxStorage) error
type CtxStorage map[string]interface{}

func NewTask(name string, id int8, execute TaskFunc, rollback TaskFunc) AtomTask {

	if name == "" || execute == nil || rollback == nil {
		fmt.Errorf("new task param failed")
		return nil
	}

	return &atomTaskBase{name: name, tId: id, exec: execute, rb: rollback, errors: &InstallerError{}}
}

type atomTaskBase struct {
	name   string
	tId    int8
	exec   TaskFunc
	rb     TaskFunc
	errors *InstallerError
}

func (a *atomTaskBase) Name() string {
	return a.name
}

func (a *atomTaskBase) Id() int8 {
	return a.tId
}

func (a *atomTaskBase) Execute(cts CtxStorage) error {
	return a.exec(cts)
}

func (a *atomTaskBase) Rollback(cts CtxStorage) error { // TODO@yiccai
	return a.rb(cts)
}

func (a *atomTaskBase) SetErrors(isExecute bool, f, t int8, fName, tName string, err error) {

	if isExecute {
		a.errors.SetExecuteErr(f, t, fName, tName, err)
	} else {
		a.errors.SetRollbackErr(f, t, fName, tName, err)
	}
}

func (a *atomTaskBase) GetErrors() (*ErrExecute, *ErrRollback) {
	return a.errors.GetExecuteErr(), a.errors.GetRollbackErr()
}
