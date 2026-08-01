package download

import (
	"io"
)

type FileInfo struct {
	Filename    string
	Type        string
	Ext         string
	ContentType string
	Size        int64
	BodyReader  io.ReadCloser
}
