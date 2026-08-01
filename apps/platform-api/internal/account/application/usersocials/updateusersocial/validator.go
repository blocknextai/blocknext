package updateusersocial

import (
	"net/url"
	"strings"

	accountApplicationUserSocials "github.com/blocknextai/platform-api/internal/account/application/usersocials"
)

func (c *UpdateUserSocialCommand) Validate() error {
	for _, item := range c.Items {
		if strings.TrimSpace(item.URL) == "" {
			return accountApplicationUserSocials.ErrURLRequired
		}

		parsedURL, err := url.Parse(item.URL)
		if err != nil {
			return accountApplicationUserSocials.ErrInvalidURL
		}

		if parsedURL.Scheme != "https" {
			return accountApplicationUserSocials.ErrHttpsRequired
		}

		if parsedURL.Host == "" {
			return accountApplicationUserSocials.ErrInvalidURL
		}
	}

	return nil
}
