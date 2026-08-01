package result

// Metadata carries optional response metadata such as pagination details.
type Metadata struct {
	Pagination *Pagination `json:"pagination,omitempty"`
}
