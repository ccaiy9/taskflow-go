package main

import (
	"fmt"
	"github.com/gogf/gf/v2/test/gtest"
	"taskflow-go/base"
	__precheck_flow "taskflow-go/flow_test/1_precheck_flow"
	__do_business_flow "taskflow-go/flow_test/2_do_business_flow"
	"testing"
)

func tmpExec(cts base.CtxStorage) error {
	fmt.Println("---------- flow-0-task-1 tmpExec Execute ------------------")
	return nil
}

func tmpRoll(cts base.CtxStorage) error {
	fmt.Println("---------- flow-0-task-1 tmpRoll rollback ------------------")
	return nil
}

func loadFlow(ctxStorage map[string]interface{}) *base.Flow {
	//g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName("../install.yaml")

	mf := base.NewFlow("main-flow", 0)
	tasktmp := base.NewTask("tmp-task-2", 1, tmpExec, tmpRoll)

	mf.SubmitTasks(
		tasktmp,
		__precheck_flow.RegisterPreCheckFlow(ctxStorage),
		__do_business_flow.RegisterPreCheckDoBusinessFlow(ctxStorage))

	return mf
}

func TestExampleExecuteLineFlow(t *testing.T) {

	gtest.C(t, func(xt *gtest.T) {

		ctxStorage := make(map[string]interface{})
		mf := loadFlow(ctxStorage)
		mf.Execute(ctxStorage)
		fmt.Println(mf.PrintErrors())

	})
}

func TestExampleManuRollBackLineFlow(t *testing.T) {

	gtest.C(t, func(xt *gtest.T) {
		ctxStorage := make(map[string]interface{})
		mf := loadFlow(ctxStorage)

		fid, tid, _ := mf.GetFailedHint()
		mf.UpdateFailedHint(fid, tid)
		mf.RollBackByManual(ctxStorage)
		fmt.Println(mf.PrintErrors())
	})
}
