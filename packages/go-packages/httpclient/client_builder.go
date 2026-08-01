package httpclient

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/base64"
)

// ClientBuilder configures and constructs a Client using a fluent API.
type ClientBuilder struct {
	context          context.Context
	baseURL          string
	timeout          time.Duration
	queryParams      map[string]string
	headers          map[string]string
	contentType      ContentType
	oauth1           *OAuth1Config
	retryConfig      *RetryConfig
	checkRedirect    func(req *http.Request, via []*http.Request) error
	transport        http.RoundTripper
	allowDestination DestinationGuard
}

// NewClientBuilder returns a ClientBuilder initialized with a background
// context and a default retry policy.
func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		context:     context.Background(),
		queryParams: make(map[string]string),
		headers:     make(map[string]string),
		retryConfig: &RetryConfig{
			MaxRetries: 3,
			BackoffMs:  1000,
		},
	}
}

// Context sets the default context used by requests built from the client.
func (c *ClientBuilder) Context(ctx context.Context) *ClientBuilder {
	c.context = ctx
	return c
}

// BaseURL sets the base URL prepended to each request path.
func (c *ClientBuilder) BaseURL(baseURL string) *ClientBuilder {
	c.baseURL = baseURL
	return c
}

// JSONContentType sets the default content type to JSON.
func (c *ClientBuilder) JSONContentType() *ClientBuilder {
	c.contentType = ContentTypeJSON
	return c
}

// ContentType sets the default content type to the given value.
func (c *ClientBuilder) ContentType(value string) *ClientBuilder {
	c.contentType = ContentType(value)
	return c
}

// FormUrlencodedContentType sets the default content type to URL-encoded form data.
func (c *ClientBuilder) FormUrlencodedContentType() *ClientBuilder {
	c.contentType = ContentTypeFormUrlencoded
	return c
}

// MultipartFormContentType sets the default content type to multipart form data.
func (c *ClientBuilder) MultipartFormContentType() *ClientBuilder {
	c.contentType = ContentTypeMultipartForm
	return c
}

// BasicAuth sets an HTTP Basic Authorization header from the username and password.
func (c *ClientBuilder) BasicAuth(username string, password string) *ClientBuilder {
	var credBuilder strings.Builder
	credBuilder.WriteString(username)
	credBuilder.WriteString(":")
	credBuilder.WriteString(password)
	credentials := credBuilder.String()

	var authBuilder strings.Builder
	authBuilder.WriteString("Basic ")
	authBuilder.WriteString(base64.Encode([]byte(credentials)))

	c.headers["Authorization"] = authBuilder.String()

	return c
}

// BearerAuth sets a Bearer token Authorization header.
func (c *ClientBuilder) BearerAuth(token string) *ClientBuilder {
	var builder strings.Builder
	builder.WriteString("Bearer ")
	builder.WriteString(token)

	c.headers["Authorization"] = builder.String()

	return c
}

// OAuth1 configures OAuth1 signing using the given consumer and access credentials.
func (c *ClientBuilder) OAuth1(consumerKey, consumerSecret, accessToken, accessTokenSecret string) *ClientBuilder {
	c.oauth1 = &OAuth1Config{
		ConsumerKey:       consumerKey,
		ConsumerSecret:    consumerSecret,
		AccessToken:       accessToken,
		AccessTokenSecret: accessTokenSecret,
	}
	return c
}

// Timeout sets the request timeout for the underlying HTTP client.
func (c *ClientBuilder) Timeout(timeout time.Duration) *ClientBuilder {
	c.timeout = timeout
	return c
}

// QueryParam adds a default query parameter applied to every request.
func (c *ClientBuilder) QueryParam(key string, value string) *ClientBuilder {
	c.queryParams[key] = value
	return c
}

// Header adds a default header applied to every request.
func (c *ClientBuilder) Header(key, value string) *ClientBuilder {
	c.headers[key] = value
	return c
}

// RetryConfig sets the maximum number of retries and the base backoff in milliseconds.
func (c *ClientBuilder) RetryConfig(maxRetries int, backoffMs int) *ClientBuilder {
	c.retryConfig = &RetryConfig{
		MaxRetries: maxRetries,
		BackoffMs:  backoffMs,
	}
	return c
}

// CheckRedirect sets the redirect policy of the underlying HTTP client.
func (c *ClientBuilder) CheckRedirect(fn func(req *http.Request, via []*http.Request) error) *ClientBuilder {
	c.checkRedirect = fn
	return c
}

// NoRedirect configures the client to not follow redirects, returning the last response instead.
func (c *ClientBuilder) NoRedirect() *ClientBuilder {
	c.checkRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return c
}

// Transport sets the HTTP round tripper used by the underlying client.
func (c *ClientBuilder) Transport(transport http.RoundTripper) *ClientBuilder {
	c.transport = transport
	return c
}

// AllowDestination sets a guard that is evaluated against every request URL
// before it is sent, and against each redirect hop unless a custom redirect
// policy is configured. Returning a non-nil error from the guard blocks the
// request. Use it to allow-list schemes and hosts and mitigate SSRF when
// request destinations may be influenced by untrusted input. For robust
// protection against DNS rebinding, combine this with a custom Transport whose
// DialContext validates the resolved IP.
func (c *ClientBuilder) AllowDestination(guard DestinationGuard) *ClientBuilder {
	c.allowDestination = guard
	return c
}

// Build constructs a Client from the builder's accumulated configuration.
func (c *ClientBuilder) Build() *Client {
	checkRedirect := c.checkRedirect
	if c.allowDestination != nil && checkRedirect == nil {
		guard := c.allowDestination
		checkRedirect = func(req *http.Request, _ []*http.Request) error {
			return guard(req.URL)
		}
	}

	return &Client{
		context: c.context,
		client: &http.Client{
			Timeout:       c.timeout,
			CheckRedirect: checkRedirect,
			Transport:     c.transport,
		},
		baseURL:          c.baseURL,
		timeout:          c.timeout,
		queryParams:      c.queryParams,
		headers:          c.headers,
		contentType:      c.contentType,
		oauth1:           c.oauth1,
		retryConfig:      c.retryConfig,
		allowDestination: c.allowDestination,
	}
}
