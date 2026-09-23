package nodes

type NodeAnnotations struct {
	ReadOnly    bool  `json:"readOnly,omitempty"`
	Destructive *bool `json:"destructive,omitempty"`
	Idempotent  bool  `json:"idempotent,omitempty"`
	OpenWorld   *bool `json:"openWorld,omitempty"`
}
