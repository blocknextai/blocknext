package apperror

import (
	"net/http"
)

// KindToStatus returns the HTTP status code corresponding to kind, defaulting
// to http.StatusInternalServerError for unknown kinds.
func KindToStatus(kind Kind) int {
	switch kind {
	case KindValidation:
		return http.StatusBadRequest
	case KindNotFound:
		return http.StatusNotFound
	case KindConflict:
		return http.StatusConflict
	case KindUnauthorized:
		return http.StatusUnauthorized
	case KindForbidden:
		return http.StatusForbidden
	case KindPayloadTooLarge:
		return http.StatusRequestEntityTooLarge
	case KindUnsupportedMedia:
		return http.StatusUnsupportedMediaType
	case KindRateLimited:
		return http.StatusTooManyRequests
	case KindBadGateway:
		return http.StatusBadGateway
	case KindUnavailable:
		return http.StatusServiceUnavailable
	case KindGatewayTimeout:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError
	}
}
