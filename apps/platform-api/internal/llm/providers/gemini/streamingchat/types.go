package streamingchat

type streamGenerateContentRequest struct {
	SystemInstruction *SystemInstruction `json:"system_instruction,omitempty"`
	Contents          []Content          `json:"contents"`
	GenerationConfig  GenerationConfig   `json:"generationConfig"`
	CachedContent     string             `json:"cachedContent,omitempty"`
}

type SystemInstruction struct {
	Parts []Part `json:"parts"`
}

type Content struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GenerationConfig struct {
	Temperature     float32 `json:"temperature"`
	TopK            int32   `json:"topK"`
	TopP            float32 `json:"topP"`
	MaxOutputTokens int32   `json:"maxOutputTokens"`
}

type SSEResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content      CandidateContent `json:"content"`
	FinishReason string           `json:"finishReason"`
}

type CandidateContent struct {
	Parts []CandidatePart `json:"parts"`
	Role  string          `json:"role"`
}

type CandidatePart struct {
	Text string `json:"text"`
}
