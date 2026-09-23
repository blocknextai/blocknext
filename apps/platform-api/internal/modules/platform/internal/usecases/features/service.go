package features

type Service struct {
	functionCallingEnabled     bool
	workflowsGenerationEnabled bool
}

func NewService(
	functionCallingEnabled bool,
	workflowsGenerationEnabled bool,
) *Service {
	return &Service{
		functionCallingEnabled:     functionCallingEnabled,
		workflowsGenerationEnabled: workflowsGenerationEnabled,
	}
}
