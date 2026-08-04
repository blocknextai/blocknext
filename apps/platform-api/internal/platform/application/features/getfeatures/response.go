package getfeatures

type GetFeaturesResponse struct {
	FunctionCalling     bool `json:"functionCalling"`
	WorkflowsGeneration bool `json:"workflowsGeneration"`
}
