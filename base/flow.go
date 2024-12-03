package base

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	p = "./installer.hint"
)

type failedHint struct {
	f int8
	t int8
}

// Flow line-flow
type Flow struct {
	name         string
	autoRollBack bool

	fId int8
	//failedMaster
	//failedSub
	failedHint

	tasks []AtomTask
}

func NewFlow(name string, id int8) *Flow {
	if name == "" {
		fmt.Errorf("new task param failed")
		return nil
	}

	return &Flow{
		name:         name,
		fId:          id,
		autoRollBack: true,
	}
}

func (f *Flow) Name() string {
	return f.name
}

func (f *Flow) Id() int8 {
	return f.fId
}

func (f *Flow) Execute(cts CtxStorage) error {

	for _, task := range f.tasks {
		if err := task.Execute(cts); err != nil {

			f.f = f.Id()
			f.t = task.Id()

			if t, ok := task.(*atomTaskBase); ok {
				t.SetErrors(true, f.Id(), t.Id(), f.Name(), t.Name(), err)
			}

			if f.autoRollBack {
				if err := f.Rollback(cts); err != nil {
					// FALLTHROUGH
					fmt.Println("Rollback failed")
					return nil
				}
			} else {
				if _, ok := task.(*atomTaskBase); ok {
					if errS := f.syncFailedHint(); errS != nil {
						fmt.Println("syncFailedHint failed")
					}
				}
			}

			return err

		}
	}

	return nil
}

func (f *Flow) Rollback(cts CtxStorage) error {

	var l int8

	if f.t == 0 {
		l = int8(len(f.tasks) - 1)
	} else {
		l = f.t - 1
	}

	if l == 0 {
		return nil
	}

	for i := l; i >= 0; i-- {

		tk := f.tasks[i]
		err := tk.Rollback(cts)
		if err != nil {
			tk.SetErrors(false, f.Id(), tk.Id(), tk.Name(), tk.Name(), err)
		}
	}

	return nil
}

func (f *Flow) RollBackByManual(cts CtxStorage) error {

	for i := f.f - 1; i >= 0; i-- {
		tk := f.tasks[i]
		err := tk.Rollback(cts)
		if err != nil {
			tk.SetErrors(false, f.Id(), tk.Id(), tk.Name(), tk.Name(), err)
		}
	}

	return nil
}

func (f *Flow) SubmitTasks(ts ...AtomTask) {
	f.tasks = append(f.tasks, ts...)
}

// SetErrors no use
func (f *Flow) SetErrors(isExecute bool, fid, tid int8, fName, tName string, err error) {
	return
}

func (f *Flow) GetErrors() (*ErrExecute, *ErrRollback) {
	return nil, nil
}

func (f *Flow) PrintErrors() string {
	var flowErrors InstallerError

	handler := func(t AtomTask) {
		execErr, robErr := t.GetErrors()
		if execErr != nil {
			flowErrors.ErrExecute.Errors = append(flowErrors.ErrExecute.Errors, execErr.Errors...)
		}
		if robErr != nil {
			flowErrors.ErrRollback.Errors = append(flowErrors.ErrRollback.Errors, robErr.Errors...)
		}
	}

	for _, task := range f.tasks {

		if at, ok := task.(*Flow); ok {
			for _, a := range at.tasks {
				handler(a)
			}
		} else {
			handler(task)
		}
	}

	reqJSON, _ := json.Marshal(flowErrors)
	return string(reqJSON)
}

func (f *Flow) UpdateFailedHint(fid, tid int8) {
	f.f = fid
	f.t = tid
}

func (f *Flow) syncFailedHint() error {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		file, err := os.Create(p)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return err
		}
		file.Close()
	}

	file, err := os.OpenFile(p, os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	s := fmt.Sprintf("%d-%d", f.failedHint.f, f.failedHint.t)
	_, err = file.WriteString(s)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	return nil
}

func (f *Flow) GetFailedHint() (int8, int8, error) {
	r, err := os.Open(p)
	if err != nil {
		fmt.Printf("open %s failed: %+v", p, err)
		return 0, 0, err
	}

	scanner := bufio.NewScanner(r)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
	}

	parts := strings.Split(line, "-")
	fid, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Println("trans error:", err)
		return 0, 0, err
	}
	tid, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println("trans error:", err)
		return 0, 0, err
	}

	return int8(fid), int8(tid), nil
}
