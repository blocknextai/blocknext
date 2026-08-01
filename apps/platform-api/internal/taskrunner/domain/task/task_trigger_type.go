package task

type TaskTriggerType string

const (
	TaskTriggerTypeManual      TaskTriggerType = "manual"
	TaskTriggerTypeRerunAll    TaskTriggerType = "rerun_all"
	TaskTriggerTypeRerunFailed TaskTriggerType = "rerun_failed"
	TaskTriggerTypeSchedule    TaskTriggerType = "schedule"
	TaskTriggerTypeWebhook     TaskTriggerType = "webhook"
	TaskTriggerTypeAPI         TaskTriggerType = "api"
)

var (
	TaskTriggerTypes = map[TaskTriggerType]struct{}{
		TaskTriggerTypeManual:      {},
		TaskTriggerTypeRerunAll:    {},
		TaskTriggerTypeRerunFailed: {},
		TaskTriggerTypeSchedule:    {},
		TaskTriggerTypeWebhook:     {},
		TaskTriggerTypeAPI:         {},
	}
)

func (t TaskTriggerType) String() string {
	return string(t)
}

func (t TaskTriggerType) IsValid() bool {
	_, ok := TaskTriggerTypes[t]
	return ok
}
