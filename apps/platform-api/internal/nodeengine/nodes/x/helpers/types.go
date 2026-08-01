package helpers

import (
	"net/http"
)

type ErrorResponse struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Type   string `json:"type"`
}

type RateLimitInfo struct {
	AppLimit24Hour           string `json:"appLimit24Hour,omitempty"`
	AppLimit24HourRemaining  string `json:"appLimit24HourRemaining,omitempty"`
	AppLimit24HourReset      string `json:"appLimit24HourReset,omitempty"`
	UserLimit24Hour          string `json:"userLimit24Hour,omitempty"`
	UserLimit24HourRemaining string `json:"userLimit24HourRemaining,omitempty"`
	UserLimit24HourReset     string `json:"userLimit24HourReset,omitempty"`
	RateLimitLimit           string `json:"rateLimitLimit,omitempty"`
	RateLimitRemaining       string `json:"rateLimitRemaining,omitempty"`
	RateLimitReset           string `json:"rateLimitReset,omitempty"`
}

func ExtractRateLimitInfo(headers http.Header) RateLimitInfo {
	return RateLimitInfo{
		AppLimit24Hour:           headers.Get("X-App-Limit-24hour-Limit"),
		AppLimit24HourRemaining:  headers.Get("X-App-Limit-24hour-Remaining"),
		AppLimit24HourReset:      headers.Get("X-App-Limit-24hour-Reset"),
		UserLimit24Hour:          headers.Get("X-User-Limit-24hour-Limit"),
		UserLimit24HourRemaining: headers.Get("X-User-Limit-24hour-Remaining"),
		UserLimit24HourReset:     headers.Get("X-User-Limit-24hour-Reset"),
		RateLimitLimit:           headers.Get("X-Rate-Limit-Limit"),
		RateLimitRemaining:       headers.Get("X-Rate-Limit-Remaining"),
		RateLimitReset:           headers.Get("X-Rate-Limit-Reset"),
	}
}

func RateLimitInfoToMap(info RateLimitInfo) map[string]any {
	return map[string]any{
		"appLimit24Hour":           info.AppLimit24Hour,
		"appLimit24HourRemaining":  info.AppLimit24HourRemaining,
		"appLimit24HourReset":      info.AppLimit24HourReset,
		"userLimit24Hour":          info.UserLimit24Hour,
		"userLimit24HourRemaining": info.UserLimit24HourRemaining,
		"userLimit24HourReset":     info.UserLimit24HourReset,
		"rateLimitLimit":           info.RateLimitLimit,
		"rateLimitRemaining":       info.RateLimitRemaining,
		"rateLimitReset":           info.RateLimitReset,
	}
}
