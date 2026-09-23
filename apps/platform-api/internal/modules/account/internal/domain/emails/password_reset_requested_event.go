package emails

type PasswordResetRequestedDomainEvent struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func (PasswordResetRequestedDomainEvent) EventName() string {
	return "account.password_reset.requested"
}
