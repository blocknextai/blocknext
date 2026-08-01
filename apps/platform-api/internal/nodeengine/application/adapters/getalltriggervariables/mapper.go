package getalltriggervariables

func MapGetAllTriggerVariablesQueryToGetAllTriggerVariablesResponse(
	variables []string,
) *GetAllTriggerVariablesResponse {
	response := make(GetAllTriggerVariablesResponse, 0, len(variables))
	response = append(response, variables...)
	return &response
}
