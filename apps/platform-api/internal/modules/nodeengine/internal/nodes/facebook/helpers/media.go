package helpers

import (
	"path/filepath"
	"strings"
)

var (
	MediaTypeMap = map[string]string{
		".mp4":  VideoType,
		".mov":  VideoType,
		".avi":  VideoType,
		".mkv":  VideoType,
		".webm": VideoType,
		".flv":  VideoType,
		".m4v":  VideoType,
		".jpg":  ImageType,
		".jpeg": ImageType,
		".png":  ImageType,
		".gif":  ImageType,
		".bmp":  ImageType,
		".webp": ImageType,
		".tiff": ImageType,
		".tif":  ImageType,
	}
)

func DetectMediaType(mediaURL string) string {
	ext := strings.ToLower(filepath.Ext(mediaURL))

	if mediaType, exists := MediaTypeMap[ext]; exists {
		return mediaType
	}

	return ImageType
}
