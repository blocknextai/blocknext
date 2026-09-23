package emails

type MagicLinkCreatedDomainEvent struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func (MagicLinkCreatedDomainEvent) EventName() string {
	return "account.magic_link.created"
}
