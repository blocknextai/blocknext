package taskexecutions

type ExecutionType string

const (
	ExecutionTypeManual   ExecutionType = "manual"
	ExecutionTypeSchedule ExecutionType = "schedule"
	ExecutionTypeWebhook  ExecutionType = "webhook"
	ExecutionTypeAPI      ExecutionType = "api"
)

var (
	ExecutionTypes = map[ExecutionType]struct{}{
		ExecutionTypeManual:   {},
		ExecutionTypeSchedule: {},
		ExecutionTypeWebhook:  {},
		ExecutionTypeAPI:      {},
	}
)

func (e ExecutionType) String() string {
	return string(e)
}

func (e ExecutionType) IsValid() bool {
	_, ok := ExecutionTypes[e]
	return ok
}
