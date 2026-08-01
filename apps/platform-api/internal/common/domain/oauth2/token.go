package oauth2

import (
	"time"
)

type Token struct {
	AccessToken  string     `json:"accessToken,omitempty"`
	RefreshToken string     `json:"refreshToken,omitempty"`
	ExpiresAt    *time.Time `json:"expiresAt,omitempty"`
}

func (t *Token) IsExpired(bufferMinutes int) bool {
	if t.ExpiresAt == nil {
		return false
	}

	buffer := time.Duration(bufferMinutes) * time.Minute
	return time.Now().UTC().Add(buffer).After(*t.ExpiresAt)
}

func (t *Token) NeedsRefresh() bool {
	return t.IsExpired(5)
}

func (t *Token) IsValid() bool {
	return t.AccessToken != ""
}

func (t *Token) HasRefreshToken() bool {
	return t.RefreshToken != ""
}

func (t *Token) SetExpiration(expiresInSeconds int) {
	if expiresInSeconds > 0 {
		expiresAt := time.Now().UTC().Add(time.Duration(expiresInSeconds) * time.Second)
		t.ExpiresAt = new(expiresAt)
	}
}
