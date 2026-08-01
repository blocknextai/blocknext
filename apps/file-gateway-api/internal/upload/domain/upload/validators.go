package upload

import (
	"slices"

	uploadDomain "github.com/blocknextai/file-gateway-api/internal/upload/domain"
)

type Validator interface {
	Validate(file *File) error
}

type MaxSizeValidator struct {
	Limit int64
}

func (v *MaxSizeValidator) Validate(file *File) error {
	if file.Size > v.Limit {
		return uploadDomain.ErrMaxSizeExceeded
	}
	return nil
}

type MIMEValidator struct {
	Allowed []string
}

func (v *MIMEValidator) Validate(file *File) error {
	if slices.Contains(v.Allowed, file.ContentType) {
		return nil
	}
	return uploadDomain.ErrMIMENotAllowed
}

func RunValidators(file *File, validators []Validator) error {
	for _, v := range validators {
		if err := v.Validate(file); err != nil {
			return err
		}
	}
	return nil
}

func CreateValidators(rule UploadRule) []Validator {
	var validators []Validator
	if rule.MaxSize > 0 {
		validators = append(validators, &MaxSizeValidator{Limit: rule.MaxSize})
	}
	if len(rule.AllowedMIMEs) > 0 {
		validators = append(validators, &MIMEValidator{Allowed: rule.AllowedMIMEs})
	}
	return validators
}
