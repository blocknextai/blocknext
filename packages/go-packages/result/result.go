// Package result provides a generic envelope for API responses, carrying a
// success flag, typed payload, optional message, metadata, and error details,
// along with helpers for constructing success and failure results.
package result

// Result is a generic response envelope wrapping a typed payload of type T
// together with status, message, metadata, and error information.
type Result[T any] struct {
	IsSuccess bool         `json:"isSuccess"`
	Data      *T           `json:"data,omitempty"`
	Message   string       `json:"message,omitempty"`
	Meta      *Metadata    `json:"meta,omitempty"`
	Error     *ErrorDetail `json:"error,omitempty"`
}

// ErrorDetail describes an error returned in a Result via a code and human-readable cause.
type ErrorDetail struct {
	Code  string `json:"code,omitempty"`
	Cause string `json:"cause,omitempty"`
}

// Ok returns a successful Result carrying val and applying any provided options.
//
// val is copied only at the top level. If T contains slices, maps or pointers,
// the returned Result shares that memory with the caller: do not mutate, append
// into or recycle val after calling Ok and before the Result is serialized.
//
// Data is always set, so the data field is always serialized for a Result built
// by Ok, including when val is the zero value of T.
func Ok[T any](val T, opts ...Option) *Result[T] {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	return &Result[T]{
		IsSuccess: true,
		Data:      &val,
		Message:   o.message,
		Meta:      o.metadata,
	}
}

// Fail returns a failed Result whose message defaults to err's text and applying
// any provided options.
//
// Fail never panics: a nil err, or an err whose Error method panics on a nil
// receiver, yields the default message instead. A WithMessage option always
// takes precedence over err's text, and an empty message falls back to the
// default so a failure is never returned without one.
//
// Only the error's own text is used. Any wrapped cause is left untouched and is
// serialized only if the caller opts in with WithCause.
func Fail[T any](err error, opts ...Option) *Result[T] {
	o := &options{message: errorMessage(err)}
	for _, opt := range opts {
		opt(o)
	}

	return &Result[T]{
		IsSuccess: false,
		Message:   o.message,
		Meta:      o.metadata,
	}
}

// WithCause sets the cause on the Result's error detail and returns the Result.
func (r *Result[T]) WithCause(cause string) *Result[T] {
	if r.Error == nil {
		r.Error = &ErrorDetail{}
	}
	r.Error.Cause = cause
	return r
}

// WithCode sets the code on the Result's error detail and returns the Result.
func (r *Result[T]) WithCode(code string) *Result[T] {
	if r.Error == nil {
		r.Error = &ErrorDetail{}
	}
	r.Error.Code = code
	return r
}

// errorMessage returns err's text, or an empty string when err is nil or its
// Error method panics on a nil receiver.
func errorMessage(err error) (message string) {
	if err == nil {
		return ""
	}

	defer func() {
		if recover() != nil {
			message = ""
		}
	}()

	return err.Error()
}
