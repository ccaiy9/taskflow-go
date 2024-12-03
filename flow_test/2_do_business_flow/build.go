package __do_business_flow

import "taskflow-go/base"

func RegisterPreCheckDoBusinessFlow(cts base.CtxStorage) *base.Flow {
	flow := base.NewFlow("RegisterPreCheckDoBusinessFlow", 2)
	flow.SubmitTasks(registerInstallToolsTask(cts), registerSendPkgTask(cts))

	return flow
}
