package functioncalling

type SuccessResponse struct {
	Candidates    []Candidate   `json:"candidates"`
	UsageMetadata UsageMetadata `json:"usageMetadata"`
}

type UsageMetadata struct {
	PromptTokenCount        int32 `json:"promptTokenCount"`
	CachedContentTokenCount int32 `json:"cachedContentTokenCount"`
	CandidatesTokenCount    int32 `json:"candidatesTokenCount"`
	ThoughtsTokenCount      int32 `json:"thoughtsTokenCount"`
	TotalTokenCount         int32 `json:"totalTokenCount"`
}

type Candidate struct {
	Content CandidateContent `json:"content"`
}

type CandidateContent struct {
	Parts []CandidatePart `json:"parts"`
}

type CandidatePart struct {
	FunctionCall FunctionCall `json:"functionCall"`
}

type FunctionCall struct {
	Name string         `json:"name"`
	Args map[string]any `json:"args"`
}
