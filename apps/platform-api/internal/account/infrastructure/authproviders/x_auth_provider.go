package authproviders

import (
	"context"
	"net/url"
	"strings"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/base64"
	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/account/application/auth/createusertoken"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	accountDomainUserNonces "github.com/blocknextai/platform-api/internal/account/domain/usernonces"
)

var (
	ErrXInvalidCode       = apperror.Validation("x invalid code")
	ErrXInvalidState      = apperror.Validation("x invalid state")
	ErrXAccessTokenFailed = apperror.Internal("x access token failed")
	ErrXUserInfoFailed    = apperror.Internal("x user info failed")
	ErrXRequestFailed     = apperror.Unavailable("x request failed")
	ErrXInvalidResponse   = apperror.Internal("x invalid response")
	ErrXInvalidConfig     = apperror.Internal("x invalid config")
)

type XAuthProvider struct {
	clientId     string
	clientSecret string
	redirectURI  string
	scope        string

	userNonceRepository accountDomainUserNonces.UserNonceRepository
}

func NewXAuthProvider(
	clientId string,
	clientSecret string,
	redirectURI string,
	userNonceRepository accountDomainUserNonces.UserNonceRepository,
) createusertoken.AuthProvider {
	return &XAuthProvider{
		clientId:            clientId,
		clientSecret:        clientSecret,
		redirectURI:         redirectURI,
		scope:               "tweet.read users.read offline.access",
		userNonceRepository: userNonceRepository,
	}
}

func (s *XAuthProvider) Validate(ctx context.Context, request createusertoken.Request) (*createusertoken.Response, error) {
	code, ok := request.Payload["code"].(string)
	if !ok {
		return nil, ErrXInvalidCode
	}
	state, ok := request.Payload["state"].(string)
	if !ok {
		return nil, ErrXInvalidState
	}

	userNonce, err := s.userNonceRepository.GetByNonceAndAuthProvider(ctx, state, accountDomain.AuthProviderX)
	if err != nil {
		return nil, err
	}

	tokenResponse, err := s.getAccessToken(code, userNonce)
	if err != nil {
		return nil, err
	}

	xUser, err := s.getXUserInfo(tokenResponse.AccessToken)
	if err != nil {
		return nil, err
	}

	return &createusertoken.Response{
		ProviderID:  xUser.Data.ID,
		Identifier:  xUser.Data.Username,
		DisplayName: xUser.Data.Name,
		Nonce:       userNonce,
	}, nil
}

func (s *XAuthProvider) BuildLoginMessage(nonce string) string {
	return ""
}

func (s *XAuthProvider) GenerateOAuthURL(userNonce *accountDomainUserNonces.UserNonce) (string, error) {
	if strings.TrimSpace(s.clientId) == "" || strings.TrimSpace(s.redirectURI) == "" {
		return "", ErrXInvalidConfig
	}

	params := url.Values{}
	params.Set("access_type", "offline")
	params.Set("response_type", "code")
	params.Set("prompt", "consent")
	params.Set("scope", s.scope)
	params.Set("client_id", s.clientId)
	params.Set("redirect_uri", s.redirectURI)
	params.Set("state", userNonce.Nonce)
	params.Set("code_challenge", userNonce.CodeChallenge)
	params.Set("code_challenge_method", userNonce.CodeChallengeMethod)

	var builder strings.Builder
	builder.WriteString("https://x.com/i/oauth2/authorize?")
	builder.WriteString(params.Encode())
	return builder.String(), nil
}

type XTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func (s *XAuthProvider) getAccessToken(code string, userNonce *accountDomainUserNonces.UserNonce) (*XTokenResponse, error) {
	client := httpclient.NewClientBuilder().
		BaseURL("https://api.x.com").
		FormUrlencodedContentType().
		Build()

	var credBuilder strings.Builder
	credBuilder.WriteString(s.clientId)
	credBuilder.WriteString(":")
	credBuilder.WriteString(s.clientSecret)

	var authBuilder strings.Builder
	authBuilder.WriteString("Basic ")
	authBuilder.WriteString(base64.Encode([]byte(credBuilder.String())))
	authorization := authBuilder.String()

	var token XTokenResponse
	response, err := client.Post("/2/oauth2/token").
		Body(map[string]string{
			"grant_type":    "authorization_code",
			"code":          code,
			"redirect_uri":  s.redirectURI,
			"scope":         s.scope,
			"code_verifier": userNonce.CodeVerifier,
		}).
		Header("Authorization", authorization).
		Do(&token, nil)

	if err != nil {
		return nil, ErrXRequestFailed
	}

	if !response.IsSuccess() {
		return nil, ErrXAccessTokenFailed
	}

	if token.AccessToken == "" {
		return nil, ErrXInvalidResponse
	}

	return &token, nil
}

type XUserInfo struct {
	Data struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Name     string `json:"name"`
	} `json:"data"`
}

func (s *XAuthProvider) getXUserInfo(accessToken string) (*XUserInfo, error) {
	client := httpclient.NewClientBuilder().
		BaseURL("https://api.x.com").
		BearerAuth(accessToken).
		Build()

	var userInfo XUserInfo
	response, err := client.Get("/2/users/me").
		QueryParam("user.fields", "id,username,name").
		Do(&userInfo, nil)

	if err != nil {
		return nil, ErrXRequestFailed
	}

	if !response.IsSuccess() {
		return nil, ErrXUserInfoFailed
	}

	if userInfo.Data.ID == "" || userInfo.Data.Username == "" {
		return nil, ErrXInvalidResponse
	}

	return &userInfo, nil
}
