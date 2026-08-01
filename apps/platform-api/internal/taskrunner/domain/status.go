package domain

type Status string

const (
	StatusPending   Status = "pending"
	StatusScheduled Status = "scheduled"
	StatusRunning   Status = "running"
	StatusRetrying  Status = "retrying"

	StatusSuccess   Status = "success"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
	StatusSkipped   Status = "skipped"
)

var (
	AllStatuses = map[Status]struct{}{
		StatusPending:   {},
		StatusScheduled: {},
		StatusRunning:   {},
		StatusRetrying:  {},
		StatusSuccess:   {},
		StatusFailed:    {},
		StatusCancelled: {},
		StatusSkipped:   {},
	}
)

func (s Status) String() string {
	return string(s)
}

func (s Status) IsValid() bool {
	_, ok := AllStatuses[s]
	return ok
}
