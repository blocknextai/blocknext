package authproviders

import (
	"context"
	"net/url"
	"strings"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/cast"
	"github.com/blocknextai/go-packages/httpclient"
	"github.com/blocknextai/platform-api/internal/account/application/auth/createusertoken"
	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	accountDomainUserNonces "github.com/blocknextai/platform-api/internal/account/domain/usernonces"
)

var (
	ErrGithubInvalidCode         = apperror.Validation("github invalid code")
	ErrGithubInvalidState        = apperror.Validation("github invalid state")
	ErrGithubTokenExchangeFailed = apperror.Internal("github token exchange failed")
	ErrGithubUserInfoFailed      = apperror.Internal("github user info failed")
	ErrGithubEmailFailed         = apperror.Internal("github email failed")
	ErrGithubMissingAccessToken  = apperror.Internal("github missing access token")
	ErrGithubMissingUserID       = apperror.Internal("github missing user id")
	ErrGithubNoVerifiedEmail     = apperror.Validation("github no verified email")
	ErrGithubInvalidConfig       = apperror.Internal("github invalid config")
)

type GithubAuthProvider struct {
	clientId     string
	clientSecret string
	redirectURI  string
	scope        string

	userNonceRepository accountDomainUserNonces.UserNonceRepository
}

func NewGithubAuthProvider(
	clientId string,
	clientSecret string,
	redirectURI string,
	userNonceRepository accountDomainUserNonces.UserNonceRepository,
) createusertoken.AuthProvider {
	return &GithubAuthProvider{
		clientId:            clientId,
		clientSecret:        clientSecret,
		redirectURI:         redirectURI,
		scope:               "read:user user:email",
		userNonceRepository: userNonceRepository,
	}
}

func (g *GithubAuthProvider) Validate(ctx context.Context, request createusertoken.Request) (*createusertoken.Response, error) {
	code, ok := request.Payload["code"].(string)
	if !ok {
		return nil, ErrGithubInvalidCode
	}
	state, ok := request.Payload["state"].(string)
	if !ok {
		return nil, ErrGithubInvalidState
	}

	nonce, err := g.userNonceRepository.GetByNonceAndAuthProvider(ctx, state, accountDomain.AuthProviderGithub)
	if err != nil {
		return nil, err
	}

	tokenResp, err := g.exchangeCodeForToken(ctx, code)
	if err != nil {
		return nil, err
	}

	userInfo, err := g.getUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(userInfo.Email) == "" {
		userInfo.Email, _ = g.getPrimaryEmail(ctx, tokenResp.AccessToken)
	}

	return &createusertoken.Response{
		ProviderID:  cast.ToString(userInfo.ID),
		Identifier:  userInfo.Login,
		DisplayName: userInfo.Name,
		Nonce:       nonce,
	}, nil
}

func (g *GithubAuthProvider) BuildLoginMessage(nonce string) string {
	return ""
}

func (g *GithubAuthProvider) GenerateOAuthURL(userNonce *accountDomainUserNonces.UserNonce) (string, error) {
	if strings.TrimSpace(g.clientId) == "" || strings.TrimSpace(g.redirectURI) == "" {
		return "", ErrGithubInvalidConfig
	}

	params := url.Values{}
	params.Set("access_type", "offline")
	params.Set("response_type", "code")
	params.Set("prompt", "consent")
	params.Set("scope", g.scope)
	params.Set("client_id", g.clientId)
	params.Set("redirect_uri", g.redirectURI)
	params.Set("state", userNonce.Nonce)

	var builder strings.Builder
	builder.WriteString("https://github.com/login/oauth/authorize?")
	builder.WriteString(params.Encode())
	return builder.String(), nil
}

type githubTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func (g *GithubAuthProvider) exchangeCodeForToken(ctx context.Context, code string) (*githubTokenResponse, error) {
	client := httpclient.NewClientBuilder().
		BaseURL("https://github.com").
		JSONContentType().
		Header("Accept", "application/json").
		Build()

	var token githubTokenResponse
	response, err := client.Post("/login/oauth/access_token").
		Body(map[string]string{
			"client_id":     g.clientId,
			"client_secret": g.clientSecret,
			"code":          code,
			"scope":         g.scope,
		}).
		Do(&token, nil)

	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		return nil, ErrGithubTokenExchangeFailed
	}

	if strings.TrimSpace(token.AccessToken) == "" {
		return nil, ErrGithubMissingAccessToken
	}

	return &token, nil
}

type githubUserInfo struct {
	ID    float64 `json:"id"`
	Login string  `json:"login"`
	Name  string  `json:"name"`
	Email string  `json:"email"`
}

func (g *GithubAuthProvider) getUserInfo(ctx context.Context, token string) (*githubUserInfo, error) {
	client := httpclient.NewClientBuilder().
		BaseURL("https://api.github.com").
		BearerAuth(token).
		Header("Accept", "application/vnd.github+json").
		Build()

	var user githubUserInfo
	response, err := client.Get("/user").
		Do(&user, nil)

	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		return nil, ErrGithubUserInfoFailed
	}

	if user.ID == 0 {
		return nil, ErrGithubMissingUserID
	}

	return &user, nil
}

func (g *GithubAuthProvider) getPrimaryEmail(ctx context.Context, token string) (string, error) {
	client := httpclient.NewClientBuilder().
		BaseURL("https://api.github.com").
		BearerAuth(token).
		Header("Accept", "application/vnd.github+json").
		Build()

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	response, err := client.Get("/user/emails").
		Do(&emails, nil)

	if err != nil {
		return "", ErrGithubEmailFailed
	}

	if !response.IsSuccess() {
		return "", ErrGithubEmailFailed
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}

	return "", ErrGithubNoVerifiedEmail
}
