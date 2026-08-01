package dag

// NodeSettings holds per-node execution settings such as retry policy,
// timeout, and error-handling behavior.
type NodeSettings struct {
	MaxRetries      float64 `json:"maxRetries"`
	RetryDelay      float64 `json:"retryDelay"`
	Timeout         float64 `json:"timeout"`
	ContinueOnError bool    `json:"continueOnError"`
	Disabled        bool    `json:"disabled"`
}
