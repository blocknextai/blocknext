package upload

type UploadRule struct {
	ID                 string   `json:"id"`
	Title              string   `json:"title"`
	MaxSize            int64    `json:"maxSize"`
	AllowedMIMEs       []string `json:"allowedMimes"`
	DefaultFolder      string   `json:"-"`
	IsOverrideFilename bool     `json:"-"`
	IsPublic           bool     `json:"-"`
}
