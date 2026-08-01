package exchangecode

func MapToResponse(status string) *ExchangeCodeResponse {
	return &ExchangeCodeResponse{
		Status: status,
	}
}
