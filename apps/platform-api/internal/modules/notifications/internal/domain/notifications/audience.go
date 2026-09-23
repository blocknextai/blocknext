package notifications

type AudienceType string

const (
	AudienceTypeUser         AudienceType = "user"
	AudienceTypeOrganization AudienceType = "organization"
)

var audienceTypes = map[AudienceType]struct{}{
	AudienceTypeUser:         {},
	AudienceTypeOrganization: {},
}

func (t AudienceType) String() string {
	return string(t)
}

func (t AudienceType) IsValid() bool {
	_, ok := audienceTypes[t]
	return ok
}
