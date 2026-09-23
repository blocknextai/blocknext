package oauth2

import (
	"slices"
)

type ResponseType string

const (
	CodeResponseType ResponseType = "code"
)

var (
	AllResponseTypes = map[ResponseType]struct{}{
		CodeResponseType: {},
	}
)

func (rt ResponseType) String() string {
	return string(rt)
}

func (rt ResponseType) IsValid() bool {
	_, ok := AllResponseTypes[rt]
	return ok
}

func SupportedResponseTypes(values []string) []ResponseType {
	responseTypes := make([]ResponseType, 0, len(values))
	for _, value := range values {
		responseType := ResponseType(value)
		if responseType.IsValid() && !slices.Contains(responseTypes, responseType) {
			responseTypes = append(responseTypes, responseType)
		}
	}
	return responseTypes
}
