package credentials

type SourceType string

const (
	SourceTypeOwner    SourceType = "owner"
	SourceTypePlatform SourceType = "platform"
)

var (
	SourceTypes = map[SourceType]struct{}{
		SourceTypeOwner:    {},
		SourceTypePlatform: {},
	}
)

func (o SourceType) String() string {
	return string(o)
}

func (o SourceType) IsValid() bool {
	_, ok := SourceTypes[o]
	return ok
}
