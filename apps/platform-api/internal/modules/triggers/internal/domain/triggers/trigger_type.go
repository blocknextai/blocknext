package triggers

type TriggerType string

const (
	TriggerTypeSchedule TriggerType = "schedule"
	TriggerTypeWebhook  TriggerType = "webhook"
)

var (
	TriggerTypes = map[TriggerType]struct{}{
		TriggerTypeSchedule: {},
		TriggerTypeWebhook:  {},
	}
)

func (t TriggerType) String() string {
	return string(t)
}

func (t TriggerType) IsValid() bool {
	_, ok := TriggerTypes[t]
	return ok
}
