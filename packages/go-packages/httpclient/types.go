// Package httpclient provides a fluent, builder-style HTTP client with support
// for JSON, form, and multipart bodies, basic, bearer, and OAuth1
// authentication, configurable retries, and streaming responses.
package httpclient

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"
)

// DestinationGuard reports whether an outbound request to the given URL is
// permitted. It returns a non-nil error to block the request. It is used to
// mitigate SSRF by allow-listing schemes/hosts or rejecting internal targets.
type DestinationGuard func(*url.URL) error

// Method is an HTTP request method.
type Method string

const (
	// MethodGet is the HTTP GET method.
	MethodGet Method = "GET"
	// MethodPost is the HTTP POST method.
	MethodPost Method = "POST"
	// MethodPut is the HTTP PUT method.
	MethodPut Method = "PUT"
	// MethodDelete is the HTTP DELETE method.
	MethodDelete Method = "DELETE"
	// MethodPatch is the HTTP PATCH method.
	MethodPatch Method = "PATCH"
)

// ContentType is the value of an HTTP Content-Type header.
type ContentType string

const (
	// ContentTypeFormUrlencoded is the URL-encoded form content type.
	ContentTypeFormUrlencoded ContentType = "application/x-www-form-urlencoded"
	// ContentTypeJSON is the JSON content type.
	ContentTypeJSON ContentType = "application/json; charset=UTF-8"
	// ContentTypeMultipartForm is the multipart form content type.
	ContentTypeMultipartForm ContentType = "multipart/form-data"
)

// MultipartFormData holds the fields and files of a multipart form request body.
type MultipartFormData struct {
	Fields []FormField
	Files  []FormFile
}

// FormField represents a single named text field in a multipart form.
type FormField struct {
	Name  string
	Value string
}

// FormFile represents a single file in a multipart form, supplied either as raw
// data or via a reader.
type FormFile struct {
	FieldName string
	FileName  string
	Data      []byte
	Reader    io.Reader
}

// OAuth1Config holds the consumer and access credentials used to sign OAuth1 requests.
type OAuth1Config struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

// RetryConfig configures the maximum number of retries and the base backoff in
// milliseconds between attempts.
type RetryConfig struct {
	MaxRetries int
	BackoffMs  int
}

// Client is a configured HTTP client that creates request builders sharing its
// base URL, headers, authentication, and retry policy.
type Client struct {
	context          context.Context
	client           *http.Client
	baseURL          string
	timeout          time.Duration
	queryParams      map[string]string
	headers          map[string]string
	contentType      ContentType
	oauth1           *OAuth1Config
	retryConfig      *RetryConfig
	allowDestination DestinationGuard
}

// RequestBuilder accumulates the method, path, query parameters, headers, body,
// and content type for a single request and executes it.
type RequestBuilder struct {
	client        *Client
	context       context.Context
	method        Method
	path          string
	queryParams   map[string]string
	headers       map[string]string
	contentType   ContentType
	body          any
	contentLength int64
}

// Response holds the result of an HTTP request, including the status code,
// buffered body, headers, and an optional streaming body reader.
type Response struct {
	Status     int
	Body       []byte
	Headers    http.Header
	BodyReader io.ReadCloser
}

// IsSuccess reports whether the response status is in the 2xx range.
func (r *Response) IsSuccess() bool {
	return checkStatus(r.Status, 200, 300)
}

// IsError reports whether the response status is in the 4xx range.
func (r *Response) IsError() bool {
	return checkStatus(r.Status, 400, 500)
}

// IsServerError reports whether the response status is 500 or greater.
func (r *Response) IsServerError() bool {
	return checkStatus(r.Status, 500, 0)
}

func checkStatus(status int, left int, right int) bool {
	if right == 0 {
		return status >= left
	}

	return status >= left && status < right
}
