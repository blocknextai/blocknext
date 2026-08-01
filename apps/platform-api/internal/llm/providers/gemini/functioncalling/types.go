package functioncalling

type SuccessResponse struct {
	Candidates []Candidate `json:"candidates"`
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
