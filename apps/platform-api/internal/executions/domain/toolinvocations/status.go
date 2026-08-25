package toolinvocations

type Status string

const (
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
)

var (
	Statuses = map[Status]struct{}{
		StatusSuccess: {},
		StatusFailed:  {},
	}
)

func (s Status) String() string {
	return string(s)
}

func (s Status) IsValid() bool {
	_, ok := Statuses[s]
	return ok
}
