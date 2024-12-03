package __do_business_flow

import (
	"fmt"
	"taskflow-go/base"
)

func registerInstallToolsTask(cts base.CtxStorage) base.AtomTask {
	return base.NewTask("registerInstallToolsTask", 1, installToolsExecute, installToolsRollback)
}

func installToolsExecute(cts base.CtxStorage) error {
	fmt.Println("---------- flow-2-task-1 -installToolsExecute ------------------")
	return fmt.Errorf("flow-2-task-1 execute failed")
}

func installToolsRollback(cts base.CtxStorage) error {
	fmt.Println("---------- flow-2-task-1 installToolsRollback ------------------")
	return nil
}
