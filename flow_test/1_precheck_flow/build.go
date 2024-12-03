package __precheck_flow

import "taskflow-go/base"

func RegisterPreCheckFlow(cts base.CtxStorage) *base.Flow {
	flow := base.NewFlow("PreCheckFlow", 1)
	flow.SubmitTasks(registerPreCheckConnectTask(cts), registerPreCheckDiskTask(cts))

	return flow
}
