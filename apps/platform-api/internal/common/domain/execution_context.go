package domain

type ExecutionContext string

const (
	ExecutionContextWorkflow ExecutionContext = "workflow"
)

var (
	ExecutionContexts = map[ExecutionContext]struct{}{
		ExecutionContextWorkflow: {},
	}
)

func (ec ExecutionContext) String() string {
	return string(ec)
}

func (ec ExecutionContext) IsValid() bool {
	_, ok := ExecutionContexts[ec]
	return ok
}
