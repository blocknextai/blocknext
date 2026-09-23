package toolinvocations

type Source string

const (
	SourceMCP      Source = "mcp"
	SourcePlatform Source = "platform"
)

var (
	Sources = map[Source]struct{}{
		SourceMCP:      {},
		SourcePlatform: {},
	}
)

func (s Source) String() string {
	return string(s)
}

func (s Source) IsValid() bool {
	_, ok := Sources[s]
	return ok
}
