package emails

type RegistrationExistingEmailNotifiedDomainEvent struct {
	Email string `json:"email"`
}

func (RegistrationExistingEmailNotifiedDomainEvent) EventName() string {
	return "account.registration.existing_email_notified"
}
