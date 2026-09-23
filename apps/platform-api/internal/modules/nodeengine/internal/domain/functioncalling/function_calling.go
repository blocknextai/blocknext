package functioncalling

type FunctionCalling struct {
	ID                   string           `json:"id,omitempty"`
	Disabled             bool             `json:"disabled,omitempty"`
	FunctionDeclarations []map[string]any `json:"functionDeclarations,omitempty"`
}

type FunctionCallingManager interface {
	GetID() string
	GetDisabled() bool
	GetFunctionDeclarations() []map[string]any
}

func (f *FunctionCalling) GetID() string {
	return f.ID
}

func (f *FunctionCalling) GetDisabled() bool {
	return f.Disabled
}

func (f *FunctionCalling) GetFunctionDeclarations() []map[string]any {
	return f.FunctionDeclarations
}
