package outboxmessages

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusProcessed  Status = "processed"
	StatusDead       Status = "dead"
)

var (
	Statuses = map[Status]struct{}{
		StatusPending:    {},
		StatusProcessing: {},
		StatusProcessed:  {},
		StatusDead:       {},
	}
)

func (s Status) String() string {
	return string(s)
}

func (s Status) IsValid() bool {
	_, ok := Statuses[s]
	return ok
}
