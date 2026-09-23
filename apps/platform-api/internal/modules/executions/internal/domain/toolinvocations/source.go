package toolinvocations

type Source string

const (
	SourceMCP Source = "mcp"
)

var (
	Sources = map[Source]struct{}{
		SourceMCP: {},
	}
)

func (s Source) String() string {
	return string(s)
}

func (s Source) IsValid() bool {
	_, ok := Sources[s]
	return ok
}
