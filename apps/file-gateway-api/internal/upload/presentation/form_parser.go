package presentation

import (
	"bytes"
	"io"
	"log/slog"
	"path/filepath"
	"strings"

	commonDomain "github.com/blocknextai/file-gateway-api/internal/common/domain"
	uploadDomain "github.com/blocknextai/file-gateway-api/internal/upload/domain"
	uploadDomainUpload "github.com/blocknextai/file-gateway-api/internal/upload/domain/upload"
	bnfile "github.com/blocknextai/go-packages/file"
	"github.com/gofiber/fiber/v3"
)

const (
	maxFilenameLength = 255
)

func parseUploadForm(c fiber.Ctx, rule *uploadDomainUpload.UploadRule) (*uploadDomainUpload.File, error) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return nil, uploadDomain.ErrInvalidFile
	}

	filename, err := sanitizeFilename(fileHeader.Filename)
	if err != nil {
		return nil, err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, uploadDomain.ErrInvalidUpload
	}

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		if closeErr := file.Close(); closeErr != nil {
			slog.ErrorContext(c.RequestCtx(), "failed to close file", "error", closeErr)
		}
		return nil, uploadDomain.ErrFileNotRead
	}

	mimeType := bnfile.MIMEType(buffer[:n]).String()
	combined := io.MultiReader(bytes.NewReader(buffer[:n]), file)
	reader := commonDomain.NewLimitedReader(combined, rule.MaxSize, uploadDomain.ErrMaxSizeExceeded)

	return &uploadDomainUpload.File{
		ContentReader: reader,
		Closer:        file,
		Filename:      filename,
		ContentType:   mimeType,
		Size:          fileHeader.Size,
	}, nil
}

func sanitizeFilename(filename string) (string, error) {
	filename = filepath.Base(filename)

	if filename == "" || filename == "." || filename == ".." {
		return "", uploadDomain.ErrInvalidFilename
	}

	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")
	filename = strings.ReplaceAll(filename, "\x00", "")
	filename = strings.ReplaceAll(filename, "\r", "")
	filename = strings.ReplaceAll(filename, "\n", "")
	filename = strings.ReplaceAll(filename, "\t", "")
	filename = strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return -1
		}
		return r
	}, filename)

	if filename == "" {
		return "", uploadDomain.ErrInvalidFilename
	}

	if len(filename) > maxFilenameLength {
		ext := filepath.Ext(filename)
		filename = filename[:maxFilenameLength-len(ext)] + ext
	}

	return filename, nil
}
