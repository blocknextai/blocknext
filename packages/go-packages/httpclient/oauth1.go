package httpclient

import (
	"crypto/hmac"
	"crypto/sha1"
	"maps"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/base64"
	"github.com/blocknextai/go-packages/uuid"
)

func (r *RequestBuilder) generateOAuth1Header(method, fullURL string) (string, error) {
	oauth := r.client.oauth1

	nonce := uuid.NewV7().String()
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)

	params := map[string]string{
		"oauth_consumer_key":     oauth.ConsumerKey,
		"oauth_nonce":            nonce,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        timestamp,
		"oauth_token":            oauth.AccessToken,
		"oauth_version":          "1.0",
	}

	signature, err := r.generateOAuth1Signature(method, fullURL, params)
	if err != nil {
		return "", err
	}
	params["oauth_signature"] = signature

	var builder strings.Builder
	var authParts []string
	for key, value := range params {
		builder.Reset()
		builder.WriteString(key)
		builder.WriteString("=\"")
		builder.WriteString(url.QueryEscape(value))
		builder.WriteString("\"")
		authParts = append(authParts, builder.String())
	}
	sort.Strings(authParts)

	builder.Reset()
	builder.WriteString("OAuth ")
	builder.WriteString(strings.Join(authParts, ", "))
	return builder.String(), nil
}

func (r *RequestBuilder) generateOAuth1Signature(method, fullURL string, params map[string]string) (string, error) {
	oauth := r.client.oauth1

	u, err := url.Parse(fullURL)
	if err != nil {
		return "", err
	}

	var baseURLBuilder strings.Builder
	baseURLBuilder.WriteString(u.Scheme)
	baseURLBuilder.WriteString("://")
	baseURLBuilder.WriteString(u.Host)
	baseURLBuilder.WriteString(u.Path)
	baseURL := baseURLBuilder.String()

	allParams := make(map[string]string)
	maps.Copy(allParams, params)
	for key, values := range u.Query() {
		if len(values) > 0 {
			allParams[key] = values[0]
		}
	}

	var keys []string
	for k := range allParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var builder strings.Builder
	var paramParts []string
	for _, key := range keys {
		builder.Reset()
		builder.WriteString(url.QueryEscape(key))
		builder.WriteString("=")
		builder.WriteString(url.QueryEscape(allParams[key]))
		paramParts = append(paramParts, builder.String())
	}
	paramString := strings.Join(paramParts, "&")

	builder.Reset()
	builder.WriteString(url.QueryEscape(method))
	builder.WriteString("&")
	builder.WriteString(url.QueryEscape(baseURL))
	builder.WriteString("&")
	builder.WriteString(url.QueryEscape(paramString))
	signatureBaseString := builder.String()

	builder.Reset()
	builder.WriteString(url.QueryEscape(oauth.ConsumerSecret))
	builder.WriteString("&")
	builder.WriteString(url.QueryEscape(oauth.AccessTokenSecret))
	signingKey := builder.String()

	h := hmac.New(sha1.New, []byte(signingKey))
	h.Write([]byte(signatureBaseString))
	signature := base64.Encode(h.Sum(nil))

	return signature, nil
}
