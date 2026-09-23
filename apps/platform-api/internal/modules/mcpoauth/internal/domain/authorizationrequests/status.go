package authorizationrequests

type Status string

const (
	PendingStatus  Status = "pending"
	ApprovedStatus Status = "approved"
	DeniedStatus   Status = "denied"
)

var (
	AllStatuses = map[Status]struct{}{
		PendingStatus:  {},
		ApprovedStatus: {},
		DeniedStatus:   {},
	}
)

func (s Status) String() string {
	return string(s)
}

func (s Status) IsValid() bool {
	_, ok := AllStatuses[s]
	return ok
}
