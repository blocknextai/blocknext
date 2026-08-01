package result

import (
	"strings"
)

// SearchRequest holds a free-text search query.
type SearchRequest struct {
	Query string `json:"query" query:"query"`
}

// Normalize trims surrounding whitespace from the query in place and returns
// the normalized value, so it is correct both as a statement and as an
// expression.
func (s *SearchRequest) Normalize() SearchRequest {
	*s = NewSearchRequest(s.Query)
	return *s
}

// NewSearchRequest returns a SearchRequest with the query trimmed of surrounding
// whitespace. The query is untrusted user input and is not escaped or otherwise
// sanitized: callers must parameterize or escape it for their own sink.
func NewSearchRequest(query string) SearchRequest {
	return SearchRequest{
		Query: strings.TrimSpace(query),
	}
}
