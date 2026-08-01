package result

type options struct {
	message  string
	metadata *Metadata
}

// Option configures optional fields of a Result during construction.
type Option func(*options)

// WithMessage returns an Option that sets the Result's message.
func WithMessage(message string) Option {
	return func(o *options) {
		o.message = message
	}
}

// WithPagination returns an Option that attaches pagination details to the
// Result's metadata. Every Result receives its own copy, so a single Option
// value may safely be reused across multiple results.
func WithPagination(pagination Pagination) Option {
	return func(o *options) {
		if o.metadata == nil {
			o.metadata = &Metadata{}
		}
		o.metadata.Pagination = new(pagination)
	}
}
