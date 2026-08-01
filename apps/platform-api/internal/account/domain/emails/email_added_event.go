package emails

type EmailAddedDomainEvent struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func (EmailAddedDomainEvent) EventName() string {
	return "account.email.added"
}
