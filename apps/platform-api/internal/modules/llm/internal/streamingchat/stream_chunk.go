package streamingchat

type Chunk struct {
	Type    ChunkType
	Content string
	Error   error
}
