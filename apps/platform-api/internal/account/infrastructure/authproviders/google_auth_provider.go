package authproviders

import (
	"context"
	"net/url"
	"strings"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/account/application/auth/createusertoken"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	accountDomainUserNonces "github.com/blocknextai/platform-api/internal/account/domain/usernonces"
)

var (
	ErrGoogleInvalidCode       = apperror.Validation("google invalid code")
	ErrGoogleInvalidState      = apperror.Validation("google invalid state")
	ErrGoogleAccessTokenFailed = apperror.Internal("google access token failed")
	ErrGoogleUserInfoFailed    = apperror.Internal("google user info failed")
	ErrGoogleInvalidConfig     = apperror.Internal("google invalid config")
	ErrGoogleEmailNotVerified  = apperror.Validation("google email not verified")
)

type GoogleAuthProvider struct {
	clientId     string
	clientSecret string
	redirectURI  string
	scope        string

	userNonceRepository accountDomainUserNonces.UserNonceRepository
}

func NewGoogleAuthProvider(
	clientId string,
	clientSecret string,
	redirectURI string,
	userNonceRepository accountDomainUserNonces.UserNonceRepository,
) createusertoken.AuthProvider {
	return &GoogleAuthProvider{
		clientId:            clientId,
		clientSecret:        clientSecret,
		redirectURI:         redirectURI,
		scope:               "openid email profile",
		userNonceRepository: userNonceRepository,
	}
}

func (s *GoogleAuthProvider) Validate(ctx context.Context, request createusertoken.Request) (*createusertoken.Response, error) {
	code, ok := request.Payload["code"].(string)
	if !ok {
		return nil, ErrGoogleInvalidCode
	}
	state, ok := request.Payload["state"].(string)
	if !ok {
		return nil, ErrGoogleInvalidState
	}

	nonce, err := s.userNonceRepository.GetByNonceAndAuthProvider(ctx, state, accountDomain.AuthProviderGoogle)
	if err != nil {
		return nil, err
	}

	tokenResponse, err := s.getAccessToken(code)
	if err != nil {
		return nil, err
	}

	googleUser, err := s.getGoogleUserInfo(tokenResponse.AccessToken)
	if err != nil {
		return nil, err
	}

	if !googleUser.EmailVerified {
		return nil, ErrGoogleEmailNotVerified
	}

	return &createusertoken.Response{
		ProviderID:  googleUser.Sub,
		Identifier:  googleUser.Email,
		DisplayName: googleUser.Name,
		Nonce:       nonce,
	}, nil
}

func (s *GoogleAuthProvider) BuildLoginMessage(nonce string) string {
	return ""
}

func (s *GoogleAuthProvider) GenerateOAuthURL(userNonce *accountDomainUserNonces.UserNonce) (string, error) {
	if strings.TrimSpace(s.clientId) == "" || strings.TrimSpace(s.redirectURI) == "" {
		return "", ErrGoogleInvalidConfig
	}

	params := url.Values{}
	params.Set("access_type", "offline")
	params.Set("response_type", "code")
	params.Set("prompt", "consent")
	params.Set("scope", s.scope)
	params.Set("client_id", s.clientId)
	params.Set("redirect_uri", s.redirectURI)
	params.Set("state", userNonce.Nonce)

	var builder strings.Builder
	builder.WriteString("https://accounts.google.com/o/oauth2/v2/auth?")
	builder.WriteString(params.Encode())
	return builder.String(), nil
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (s *GoogleAuthProvider) getAccessToken(code string) (*TokenResponse, error) {
	client := httpclient.NewClientBuilder().
		BaseURL("https://oauth2.googleapis.com").
		FormUrlencodedContentType().
		Build()

	var token TokenResponse
	response, err := client.Post("/token").
		Body(map[string]string{
			"code":          code,
			"client_id":     s.clientId,
			"client_secret": s.clientSecret,
			"redirect_uri":  s.redirectURI,
			"scope":         s.scope,
			"grant_type":    "authorization_code",
		}).
		Do(&token, nil)

	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		return nil, ErrGoogleAccessTokenFailed
	}

	return &token, nil
}

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func (s *GoogleAuthProvider) getGoogleUserInfo(accessToken string) (*GoogleUserInfo, error) {
	client := httpclient.NewClientBuilder().
		BaseURL("https://openidconnect.googleapis.com/v1").
		BearerAuth(accessToken).
		Build()

	var userInfo GoogleUserInfo
	response, err := client.Get("/userinfo").
		Do(&userInfo, nil)

	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		return nil, ErrGoogleUserInfoFailed
	}

	return &userInfo, nil
}
