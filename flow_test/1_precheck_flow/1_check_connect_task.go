package __precheck_flow

import (
	"fmt"
	"taskflow-go/base"
)

func registerPreCheckConnectTask(cts base.CtxStorage) base.AtomTask {
	return base.NewTask("PreCheckConnectTask", 1, checkConnectExecute, checkConnectRollback)
}

func checkConnectExecute(cts base.CtxStorage) error {
	fmt.Println("---------- flow-1-task-1 checkConnectExecute ------------------")
	return nil
}

func checkConnectRollback(cts base.CtxStorage) error {
	fmt.Println("---------- flow-1-task-1 checkConnectRollback ------------------")
	return fmt.Errorf("flow-1-task-1 checkConnectRollback failed")
}
