package streamingchat

import (
	"context"
)

type Provider interface {
	StreamChat(ctx context.Context, systemInstruction string, messages []Message) (<-chan Chunk, error)
}
