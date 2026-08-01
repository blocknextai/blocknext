package adapters

type Adapter struct {
	ID       string `json:"id,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

type AdapterManager interface {
	GetID() string
	GetDisabled() bool
	Adapt(raw map[string]any) (*TriggerContext, error)
}

type TriggerAdapter = AdapterManager

func (a *Adapter) GetID() string {
	return a.ID
}

func (a *Adapter) GetDisabled() bool {
	return a.Disabled
}
