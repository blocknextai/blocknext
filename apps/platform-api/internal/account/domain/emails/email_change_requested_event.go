package emails

type EmailChangeRequestedDomainEvent struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func (EmailChangeRequestedDomainEvent) EventName() string {
	return "account.email_change.requested"
}
