package emails

type EmailVerificationRequestedDomainEvent struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func (EmailVerificationRequestedDomainEvent) EventName() string {
	return "account.email_verification.requested"
}
