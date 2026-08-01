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
	ErrFacebookInvalidCode       = apperror.Validation("facebook invalid code")
	ErrFacebookInvalidState      = apperror.Validation("facebook invalid state")
	ErrFacebookAccessTokenFailed = apperror.Internal("facebook access token failed")
	ErrFacebookUserInfoFailed    = apperror.Internal("facebook user info failed")
	ErrFacebookInvalidConfig     = apperror.Internal("facebook invalid config")
)

type FacebookAuthProvider struct {
	clientId     string
	clientSecret string
	redirectURI  string
	scope        string

	userNonceRepository accountDomainUserNonces.UserNonceRepository
}

func NewFacebookAuthProvider(
	clientId string,
	clientSecret string,
	redirectURI string,
	userNonceRepository accountDomainUserNonces.UserNonceRepository,
) createusertoken.AuthProvider {
	return &FacebookAuthProvider{
		clientId:            clientId,
		clientSecret:        clientSecret,
		redirectURI:         redirectURI,
		scope:               "email public_profile",
		userNonceRepository: userNonceRepository,
	}
}

func (s *FacebookAuthProvider) Validate(ctx context.Context, request createusertoken.Request) (*createusertoken.Response, error) {
	code, ok := request.Payload["code"].(string)
	if !ok {
		return nil, ErrFacebookInvalidCode
	}
	state, ok := request.Payload["state"].(string)
	if !ok {
		return nil, ErrFacebookInvalidState
	}

	nonce, err := s.userNonceRepository.GetByNonceAndAuthProvider(ctx, state, accountDomain.AuthProviderFacebook)
	if err != nil {
		return nil, err
	}

	tokenResponse, err := s.getAccessToken(code)
	if err != nil {
		return nil, err
	}

	facebookUser, err := s.getFacebookUserInfo(tokenResponse.AccessToken)
	if err != nil {
		return nil, err
	}

	return &createusertoken.Response{
		ProviderID:  facebookUser.ID,
		Identifier:  facebookUser.Email,
		DisplayName: facebookUser.Name,
		Nonce:       nonce,
	}, nil
}

func (s *FacebookAuthProvider) BuildLoginMessage(nonce string) string {
	return ""
}

func (s *FacebookAuthProvider) GenerateOAuthURL(userNonce *accountDomainUserNonces.UserNonce) (string, error) {
	if strings.TrimSpace(s.clientId) == "" || strings.TrimSpace(s.redirectURI) == "" {
		return "", ErrFacebookInvalidConfig
	}

	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("scope", s.scope)
	params.Set("client_id", s.clientId)
	params.Set("redirect_uri", s.redirectURI)
	params.Set("state", userNonce.Nonce)

	var builder strings.Builder
	builder.WriteString("https://www.facebook.com/v23.0/dialog/oauth?")
	builder.WriteString(params.Encode())
	return builder.String(), nil
}

type FacebookTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (s *FacebookAuthProvider) getAccessToken(code string) (*FacebookTokenResponse, error) {
	client := httpclient.NewClientBuilder().
		BaseURL("https://graph.facebook.com/v23.0").
		FormUrlencodedContentType().
		Build()

	var token FacebookTokenResponse
	response, err := client.Post("/oauth/access_token").
		Body(map[string]string{
			"code":          code,
			"client_id":     s.clientId,
			"client_secret": s.clientSecret,
			"redirect_uri":  s.redirectURI,
			"grant_type":    "authorization_code",
		}).
		Do(&token, nil)

	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		return nil, ErrFacebookAccessTokenFailed
	}

	return &token, nil
}

type FacebookUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (s *FacebookAuthProvider) getFacebookUserInfo(accessToken string) (*FacebookUserInfo, error) {
	client := httpclient.NewClientBuilder().
		BaseURL("https://graph.facebook.com/v23.0").
		QueryParam("access_token", accessToken).
		Build()

	var userInfo FacebookUserInfo
	response, err := client.Get("/me").
		QueryParam("fields", "id,email,name").
		Do(&userInfo, nil)

	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		return nil, ErrFacebookUserInfoFailed
	}

	return &userInfo, nil
}
