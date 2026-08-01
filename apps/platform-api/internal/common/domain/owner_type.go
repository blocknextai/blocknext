package domain

type OwnerType string

const (
	OwnerTypeOrganization OwnerType = "organization"
	OwnerTypeUser         OwnerType = "user"
)

var (
	OwnerTypes = map[OwnerType]struct{}{
		OwnerTypeOrganization: {},
		OwnerTypeUser:         {},
	}
)

func (o OwnerType) String() string {
	return string(o)
}

func (o OwnerType) IsValid() bool {
	_, ok := OwnerTypes[o]
	return ok
}
