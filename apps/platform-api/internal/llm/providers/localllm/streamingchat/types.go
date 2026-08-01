package streamingchat

type chatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float32       `json:"temperature"`
	TopK        int32         `json:"top_k,omitempty"`
	TopP        float32       `json:"top_p"`
	MaxTokens   int32         `json:"max_tokens,omitempty"`
	Stream      bool          `json:"stream"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type SSEResponse struct {
	Choices []SSEChoice `json:"choices"`
	Error   *SSEError   `json:"error"`
}

type SSEError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    any    `json:"code"`
}

type SSEChoice struct {
	Delta        SSEDelta `json:"delta"`
	FinishReason string   `json:"finish_reason"`
}

type SSEDelta struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
