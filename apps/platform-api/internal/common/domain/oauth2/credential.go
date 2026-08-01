package oauth2

type Credential struct {
	ClientID     string `json:"clientId,omitempty"`
	ClientKey    string `json:"clientKey,omitempty"`
	ClientSecret string `json:"clientSecret,omitempty"`
	OAuthToken   Token  `json:"oauthToken"`
}

func (c *Credential) GetClientIdentifier() string {
	if c.ClientID != "" {
		return c.ClientID
	}
	return c.ClientKey
}

func (c *Credential) HasClientCredentials() bool {
	clientID := c.GetClientIdentifier()
	return clientID != "" && c.ClientSecret != ""
}
