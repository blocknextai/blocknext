// Package file provides helpers for detecting MIME types, classifying files by
// extension, and reading file content.
package file

import (
	"io"
	"net/url"
	"path"

	"github.com/gabriel-vasile/mimetype"
)

// MIMEType returns the detected MIME type of the given content.
func MIMEType(content []byte) *mimetype.MIME {
	return mimetype.Detect(content)
}

// ExtensionFromMIMEType returns the file extension for the given MIME type, or
// ".bin" when the type is unknown.
func ExtensionFromMIMEType(mimeType string) string {
	mime := mimetype.Lookup(mimeType)
	if mime != nil {
		return mime.Extension()
	}
	return ".bin"
}

// ExtractFilename returns the base filename of the given URL's path, or
// "unknown" when the URL cannot be parsed.
func ExtractFilename(rawurl string) string {
	parsed, err := url.Parse(rawurl)
	if err != nil {
		return "unknown"
	}
	return path.Base(parsed.Path)
}

// Category returns a broad category, such as "image" or "document", for the
// given file extension, or "other" when it is unrecognized.
func Category(ext string) string {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return "image"
	case ".mp4", ".mov", ".avi", ".mkv", ".webm":
		return "video"
	case ".mp3", ".wav", ".ogg", ".m4a", ".aac":
		return "audio"
	case ".txt", ".csv", ".json", ".xml", ".html", ".css", ".js":
		return "text"
	case ".pdf":
		return "pdf"
	case ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx":
		return "document"
	case ".zip", ".rar", ".7z", ".tar", ".gz", ".bz2", ".xz":
		return "archive"
	case ".bin", ".exe", ".dll", ".so", ".dylib", ".jar", ".war", ".ear", ".dmg", ".pkg", ".deb", ".rpm", ".msi":
		return "binary"
	}
	return "other"
}

// ReaderToBytes reads all data from reader and returns it as a byte slice.
func ReaderToBytes(reader io.Reader) ([]byte, error) {
	return io.ReadAll(reader)
}
