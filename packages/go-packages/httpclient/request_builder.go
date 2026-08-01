package httpclient

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"strings"

	"github.com/blocknextai/go-packages/json"
)

// NewRequestBuilder returns a RequestBuilder seeded with the client's default
// query parameters and a GET method.
func (c *Client) NewRequestBuilder() *RequestBuilder {
	queryParams := make(map[string]string, len(c.queryParams))
	maps.Copy(queryParams, c.queryParams)

	return &RequestBuilder{
		client:        c,
		method:        MethodGet,
		path:          "",
		queryParams:   queryParams,
		headers:       make(map[string]string),
		contentType:   ContentTypeJSON,
		body:          nil,
		contentLength: -1,
	}
}

// Get returns a RequestBuilder for a GET request to the given path.
func (c *Client) Get(path string) *RequestBuilder {
	return c.NewRequestBuilder().Method(MethodGet).Path(path)
}

// Post returns a RequestBuilder for a POST request to the given path.
func (c *Client) Post(path string) *RequestBuilder {
	return c.NewRequestBuilder().Method(MethodPost).Path(path)
}

// Put returns a RequestBuilder for a PUT request to the given path.
func (c *Client) Put(path string) *RequestBuilder {
	return c.NewRequestBuilder().Method(MethodPut).Path(path)
}

// Delete returns a RequestBuilder for a DELETE request to the given path.
func (c *Client) Delete(path string) *RequestBuilder {
	return c.NewRequestBuilder().Method(MethodDelete).Path(path)
}

// Patch returns a RequestBuilder for a PATCH request to the given path.
func (c *Client) Patch(path string) *RequestBuilder {
	return c.NewRequestBuilder().Method(MethodPatch).Path(path)
}

// Context sets the context for this request, overriding the client default.
func (r *RequestBuilder) Context(ctx context.Context) *RequestBuilder {
	r.context = ctx
	return r
}

// Method sets the HTTP method for this request.
func (r *RequestBuilder) Method(method Method) *RequestBuilder {
	r.method = method
	return r
}

// Path sets the request path appended to the client's base URL.
func (r *RequestBuilder) Path(path string) *RequestBuilder {
	r.path = path
	return r
}

// QueryParam sets a query parameter on this request.
func (r *RequestBuilder) QueryParam(key string, value string) *RequestBuilder {
	r.queryParams[key] = value
	return r
}

// Header sets a header on this request, overriding any client default with the same key.
func (r *RequestBuilder) Header(key string, value string) *RequestBuilder {
	r.headers[key] = value
	return r
}

// Body sets the request body, which is encoded according to the content type.
func (r *RequestBuilder) Body(body any) *RequestBuilder {
	r.body = body
	return r
}

// BodyReader sets the request body to be streamed from the given reader.
func (r *RequestBuilder) BodyReader(reader io.Reader) *RequestBuilder {
	r.body = reader
	return r
}

// JSONContentType sets the content type of this request to JSON.
func (r *RequestBuilder) JSONContentType() *RequestBuilder {
	r.contentType = ContentTypeJSON
	return r
}

// ContentType sets the content type of this request to the given value.
func (r *RequestBuilder) ContentType(value string) *RequestBuilder {
	r.contentType = ContentType(value)
	return r
}

// ContentLength sets the Content-Length of this request.
func (r *RequestBuilder) ContentLength(length int64) *RequestBuilder {
	r.contentLength = length
	return r
}

// FormUrlencodedContentType sets the content type of this request to URL-encoded form data.
func (r *RequestBuilder) FormUrlencodedContentType() *RequestBuilder {
	r.contentType = ContentTypeFormUrlencoded
	return r
}

// MultipartFormContentType sets the content type of this request to multipart form data.
func (r *RequestBuilder) MultipartFormContentType() *RequestBuilder {
	r.contentType = ContentTypeMultipartForm
	return r
}

// MultipartFormBody returns a MultipartFormBuilder for composing a multipart form body.
func (r *RequestBuilder) MultipartFormBody() *MultipartFormBuilder {
	return &MultipartFormBuilder{
		requestBuilder: r,
		formData:       &MultipartFormData{},
	}
}

func (r *RequestBuilder) resolveContext() context.Context {
	if r.context != nil {
		return r.context
	}
	return r.client.context
}

func (r *RequestBuilder) send(ctx context.Context) (*http.Response, error) {
	builtURL, err := r.buildURL()
	if err != nil {
		return nil, err
	}

	if guard := r.client.allowDestination; guard != nil {
		u, err := url.Parse(builtURL)
		if err != nil {
			return nil, ErrInvalidURL
		}
		if err := guard(u); err != nil {
			return nil, err
		}
	}

	req, err := r.buildRequest(ctx, builtURL)
	if err != nil {
		return nil, err
	}

	return r.client.client.Do(req)
}

func closeResponseBody(body io.Closer) {
	if err := body.Close(); err != nil {
		slog.Error("Failed to close response body",
			"component", "RequestBuilder",
			"error", err)
	}
}

func readResponse(resp *http.Response) (*Response, error) {
	defer closeResponseBody(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Status:  resp.StatusCode,
		Body:    bodyBytes,
		Headers: resp.Header,
	}, nil
}

// Do executes the request with retries, decoding the response body into
// successTarget on success or errorTarget otherwise.
func (r *RequestBuilder) Do(successTarget any, errorTarget any) (*Response, error) {
	return r.executeWithRetry(r.resolveContext(), func() (*Response, error) {
		resp, err := r.send(r.resolveContext())
		if err != nil {
			return nil, err
		}

		response, err := readResponse(resp)
		if err != nil {
			return nil, err
		}

		target := successTarget
		if !response.IsSuccess() {
			target = errorTarget
		}

		if target != nil {
			if err := json.Unmarshal(response.Body, target); err != nil {
				slog.Error("failed to unmarshal response",
					"package", "httpclient",
					"success", response.IsSuccess(),
					"error", err,
				)
			}
		}

		return response, nil
	})
}

// DoRaw executes the request with retries and returns the response with its body
// buffered, without decoding.
func (r *RequestBuilder) DoRaw() (*Response, error) {
	return r.executeWithRetry(r.resolveContext(), func() (*Response, error) {
		resp, err := r.send(r.resolveContext())
		if err != nil {
			return nil, err
		}

		return readResponse(resp)
	})
}

// DoStream executes the request with retries and returns the response with its
// body exposed as a streaming reader that the caller must close.
func (r *RequestBuilder) DoStream() (*Response, error) {
	return r.executeWithRetry(r.resolveContext(), func() (*Response, error) {
		resp, err := r.send(r.resolveContext())
		if err != nil {
			return nil, err
		}

		return &Response{
			Status:     resp.StatusCode,
			Headers:    resp.Header,
			BodyReader: resp.Body,
		}, nil
	})
}

func (r *RequestBuilder) buildURL() (string, error) {
	var builder strings.Builder
	builder.WriteString(r.client.baseURL)
	builder.WriteString(r.path)

	u, err := url.Parse(builder.String())
	if err != nil {
		return "", ErrInvalidURL
	}

	q := u.Query()
	for k, v := range r.queryParams {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (r *RequestBuilder) buildRequest(ctx context.Context, fullURL string) (*http.Request, error) {
	var buf io.Reader
	if r.body != nil {
		switch body := r.body.(type) {
		case io.Reader:
			buf = body
		case []byte:
			buf = bytes.NewReader(body)
		case string:
			buf = strings.NewReader(body)
		default:
			switch r.contentType {
			case ContentTypeJSON:
				j, err := json.Marshal(body)
				if err != nil {
					return nil, err
				}
				buf = bytes.NewBuffer(j)
			case ContentTypeFormUrlencoded:
				form, ok := body.(url.Values)
				if !ok {
					return nil, ErrInvalidBody
				}
				buf = bytes.NewBufferString(form.Encode())
			case ContentTypeMultipartForm:
				formData, ok := body.(*MultipartFormData)
				if !ok {
					return nil, ErrInvalidBody
				}
				multipartBuf, contentType, err := r.buildMultipartForm(ctx, formData)
				if err != nil {
					return nil, err
				}
				buf = multipartBuf
				r.contentType = ContentType(contentType)
			}
		}
	}

	req, err := http.NewRequestWithContext(ctx, string(r.method), fullURL, buf)
	if err != nil {
		return nil, err
	}

	if r.contentLength >= 0 {
		req.ContentLength = r.contentLength
	}

	if r.client.oauth1 != nil {
		authHeader, err := r.generateOAuth1Header(string(r.method), fullURL)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", authHeader)
	}

	for k, v := range r.client.headers {
		req.Header.Set(k, v)
	}
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	if r.body != nil {
		if strings.TrimSpace(string(r.client.contentType)) != "" {
			req.Header.Set("Content-Type", string(r.client.contentType))
		}

		if strings.TrimSpace(string(r.contentType)) != "" {
			req.Header.Set("Content-Type", string(r.contentType))
		}
	}

	return req, nil
}
