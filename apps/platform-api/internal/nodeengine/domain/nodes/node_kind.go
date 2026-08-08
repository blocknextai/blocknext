package nodes

type NodeKind string

const (
	NodeKindAction NodeKind = "action"
	NodeKindNote   NodeKind = "note"
)

func (k NodeKind) String() string {
	return string(k)
}
