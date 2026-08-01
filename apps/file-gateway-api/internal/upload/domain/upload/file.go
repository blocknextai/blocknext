package upload

import (
	"io"
)

type File struct {
	ContentReader io.Reader
	Closer        io.Closer
	Filename      string
	ContentType   string
	Size          int64
}

func (f *File) Close() error {
	if f.Closer != nil {
		return f.Closer.Close()
	}
	return nil
}
