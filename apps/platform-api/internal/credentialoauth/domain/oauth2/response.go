package oauth2

import (
	"strings"

	commonDomainOAuth2 "github.com/blocknextai/platform-api/internal/common/domain/oauth2"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func (tr *TokenResponse) ToToken() *commonDomainOAuth2.Token {
	token := &commonDomainOAuth2.Token{
		AccessToken:  tr.AccessToken,
		RefreshToken: tr.RefreshToken,
	}

	token.SetExpiration(tr.ExpiresIn)
	return token
}

const (
	InvalidGrantError = "invalid_grant"
)

type TokenErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (ter *TokenErrorResponse) GetErrorMessage() string {
	if strings.TrimSpace(ter.Error) == "" {
		return "unknown oauth error"
	}

	if ter.ErrorDescription != "" {
		return ter.Error + ": " + ter.ErrorDescription
	}

	return ter.Error
}
