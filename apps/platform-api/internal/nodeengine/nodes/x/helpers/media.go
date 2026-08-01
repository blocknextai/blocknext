package helpers

import (
	"path/filepath"
	"strings"
)

const (
	GifCategory   = "tweet_gif"
	ImageCategory = "tweet_image"
	VideoCategory = "tweet_video"
)

const (
	ProcessingInfoStateSucceeded  = "succeeded"
	ProcessingInfoStateInProgress = "in_progress"
	ProcessingInfoStatePending    = "pending"
	ProcessingInfoStateFailed     = "failed"
)

var (
	MediaCategoryMap = map[string]string{
		".mp4":  VideoCategory,
		".mov":  VideoCategory,
		".avi":  VideoCategory,
		".mkv":  VideoCategory,
		".webm": VideoCategory,
		".flv":  VideoCategory,
		".m4v":  VideoCategory,
		".jpg":  ImageCategory,
		".jpeg": ImageCategory,
		".png":  ImageCategory,
		".bmp":  ImageCategory,
		".webp": ImageCategory,
		".tiff": ImageCategory,
		".tif":  ImageCategory,
		".gif":  GifCategory,
	}

	ContentTypeCategoryMap = map[string]string{
		"video/":    VideoCategory,
		"image/":    ImageCategory,
		"image/gif": GifCategory,
	}
)

func DetectMediaCategory(mediaURL string) string {
	ext := strings.ToLower(filepath.Ext(mediaURL))

	if mediaCategory, exists := MediaCategoryMap[ext]; exists {
		return mediaCategory
	}

	return ImageCategory
}

func DetectMediaCategoryFromContentType(contentType string, mediaURL string) string {
	if contentType != "" {
		contentType = strings.ToLower(contentType)

		for prefix, category := range ContentTypeCategoryMap {
			if strings.HasPrefix(contentType, prefix) {
				return category
			}
		}
	}

	return DetectMediaCategory(mediaURL)
}
