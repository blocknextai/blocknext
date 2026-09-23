package contract

import (
	llmStreamingchat "github.com/blocknextai/platform-api/internal/modules/llm/internal/streamingchat"
)

type Chunk = llmStreamingchat.Chunk

const ChunkTypeDone = llmStreamingchat.ChunkTypeDone
const ChunkTypeError = llmStreamingchat.ChunkTypeError
const ChunkTypeText = llmStreamingchat.ChunkTypeText

type Message = llmStreamingchat.Message
type Provider = llmStreamingchat.Provider
