package upload

import (
	uploadDomain "github.com/blocknextai/file-gateway-api/internal/upload/domain"
)

var (
	imageMIMEs = []string{
		"image/png",
		"image/jpeg",
		"image/jpg",
		"image/avif",
		"image/heic",
		"image/heif",
		"image/webp",
		"image/gif",
		"image/heif-sequence",
		"image/heic-sequence",
	}

	basicImageMIMEs = []string{
		"image/png",
		"image/jpeg",
		"image/jpg",
	}

	videoMIMEs = []string{
		"video/mp4",
		"video/webm",
		"video/ogg",
	}

	audioMIMEs = []string{
		"audio/mpeg",
		"audio/mp3",
		"audio/mp4",
		"audio/m4a",
		"audio/m4b",
		"audio/m4p",
		"audio/m4v",
	}
)

var apiElevenLabsNode = UploadRule{
	ID:                 "8044ae7e-4d97-4531-b945-94f824d71987",
	Title:              "api.nodes.elevenlabs",
	MaxSize:            10 * 1024 * 1024,
	AllowedMIMEs:       audioMIMEs,
	DefaultFolder:      "nodes/elevenlabs/",
	IsOverrideFilename: true,
	IsPublic:           true,
}

var apiGeminiNanoBananaNode = UploadRule{
	ID:                 "04ed7773-3a64-403b-a548-80a78d6ca11f",
	Title:              "api.nodes.gemini.nano_banana",
	MaxSize:            2 * 1024 * 1024,
	AllowedMIMEs:       basicImageMIMEs,
	DefaultFolder:      "nodes/gemini/nano-banana/",
	IsOverrideFilename: true,
	IsPublic:           true,
}

var apiVeoNode = UploadRule{
	ID:                 "60b9c9f6-e2e6-4782-92a0-8ba2b8c32720",
	Title:              "api.nodes.veo",
	MaxSize:            10 * 1024 * 1024,
	AllowedMIMEs:       videoMIMEs,
	DefaultFolder:      "nodes/veo/",
	IsOverrideFilename: true,
	IsPublic:           true,
}

var apiGoogleDriveGetFileNode = UploadRule{
	ID:                 "50667bea-8f4c-4255-a54e-19d32d59927d",
	Title:              "api.nodes.google.drive.get_file",
	MaxSize:            10 * 1024 * 1024,
	AllowedMIMEs:       imageMIMEs,
	DefaultFolder:      "nodes/google/drive/get-file/",
	IsOverrideFilename: true,
	IsPublic:           true,
}

var uploadRules = map[string]UploadRule{
	apiElevenLabsNode.ID:         apiElevenLabsNode,
	apiGeminiNanoBananaNode.ID:   apiGeminiNanoBananaNode,
	apiVeoNode.ID:                apiVeoNode,
	apiGoogleDriveGetFileNode.ID: apiGoogleDriveGetFileNode,
}

func GetUploadRule(uploadID string) (*UploadRule, error) {
	if uploadID == "" {
		return nil, uploadDomain.ErrInvalidUploadID
	}

	rule, exists := uploadRules[uploadID]
	if !exists {
		return nil, uploadDomain.ErrUploadRuleNotFound
	}

	return &rule, nil
}
