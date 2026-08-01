package adapters

type AdapterBuilder struct {
	adapter Adapter
}

func NewAdapterBuilder() *AdapterBuilder {
	return &AdapterBuilder{}
}

func (b *AdapterBuilder) ID(id string) *AdapterBuilder {
	b.adapter.ID = id
	return b
}

func (b *AdapterBuilder) Build() Adapter {
	return b.adapter
}
