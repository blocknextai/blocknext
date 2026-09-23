package verificationtokens

type Purpose string

const (
	PurposeEmailVerify   Purpose = "email_verify"
	PurposePasswordReset Purpose = "password_reset"
	PurposeEmailChange   Purpose = "email_change"
	PurposeMagicLink     Purpose = "magic_link"
)

var purposes = map[Purpose]struct{}{
	PurposeEmailVerify:   {},
	PurposePasswordReset: {},
	PurposeEmailChange:   {},
	PurposeMagicLink:     {},
}

func (p Purpose) String() string {
	return string(p)
}

func (p Purpose) IsValid() bool {
	_, ok := purposes[p]
	return ok
}
