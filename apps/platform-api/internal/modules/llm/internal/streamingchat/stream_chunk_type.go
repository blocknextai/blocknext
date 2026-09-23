package streamingchat

type ChunkType string

const (
	ChunkTypeText  ChunkType = "text"
	ChunkTypeDone  ChunkType = "done"
	ChunkTypeError ChunkType = "error"
)

var (
	AllChunkTypes = map[ChunkType]struct{}{
		ChunkTypeText:  {},
		ChunkTypeDone:  {},
		ChunkTypeError: {},
	}
)

func (c ChunkType) String() string {
	return string(c)
}

func (c ChunkType) IsValid() bool {
	_, ok := AllChunkTypes[c]
	return ok
}
