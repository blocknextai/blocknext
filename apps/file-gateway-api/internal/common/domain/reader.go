package domain

import (
	"io"
)

type LimitedReader struct {
	Reader       io.Reader
	Remaining    int64
	LimitReached error
}

func NewLimitedReader(reader io.Reader, limit int64, limitReachedErr error) *LimitedReader {
	return &LimitedReader{
		Reader:       reader,
		Remaining:    limit,
		LimitReached: limitReachedErr,
	}
}

func (l *LimitedReader) Read(p []byte) (int, error) {
	if l.Remaining <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > l.Remaining+1 {
		p = p[:l.Remaining+1]
	}
	n, err := l.Reader.Read(p)
	l.Remaining -= int64(n)
	if l.Remaining < 0 {
		return n, l.LimitReached
	}
	return n, err
}
