package __precheck_flow

import (
	"fmt"
	"taskflow-go/base"
)

func registerPreCheckDiskTask(cts base.CtxStorage) base.AtomTask {
	return base.NewTask("PreCheckDiskTask", 2, checkDiskExecute, checkDiskRollback)
}

func checkDiskExecute(cts base.CtxStorage) error {
	fmt.Println("---------- flow-1-task-2 checkDiskExecute ------------------")
	return nil
}

func checkDiskRollback(cts base.CtxStorage) error {
	fmt.Println("---------- flow-1-task-2 checkDiskRollback ------------------")
	return fmt.Errorf("flow-1-task-2 checkDiskRollback failed")
}
