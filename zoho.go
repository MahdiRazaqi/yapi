package zoho

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

// Config model
type Config struct {
	conf *oauth2.Config
}

// Zoho model
type Zoho struct {
	Config
	token     *oauth2.Token
	accountID string
}

// New zoho config model
func New(clientID, clientSecret, redirectURL, authURL, tokenURL string, scopes []string) *Config {
	return &Config{
		conf: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       scopes,
			RedirectURL:  redirectURL,
			Endpoint: oauth2.Endpoint{
				AuthURL:   authURL,
				TokenURL:  tokenURL,
				AuthStyle: oauth2.AuthStyleInParams,
			},
		},
	}
}

// GetAuthURI get authentication URI for recive code
func (c *Config) GetAuthURI() string {
	return fmt.Sprintf("https://accounts.zoho.com/oauth/v2/auth?scope=%v&client_id=%v&response_type=code&access_type=offline&redirect_uri=%v", strings.Join(c.conf.Scopes, ","), c.conf.ClientID, c.conf.RedirectURL)
}

// SetTokenWithCode set token from code
func (c *Config) SetTokenWithCode(code string) (*Zoho, error) {
	t, err := c.conf.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	z := &Zoho{
		Config: Config{conf: c.conf},
		token:  t,
	}

	if err := z.setAccountID(); err != nil {
		return nil, err
	}

	return z, nil
}

// SetToken set token with manual set token
func (c *Config) SetToken(t *oauth2.Token) (*Zoho, error) {
	z := &Zoho{
		Config: Config{conf: c.conf},
		token:  t,
	}

	if err := z.setAccountID(); err != nil {
		return nil, err
	}

	return z, nil
}

// RenewToken renew token
func (z *Zoho) RenewToken() error {
	if z.token.Expiry.Before(time.Now()) {
		src := z.Config.conf.TokenSource(oauth2.NoContext, z.token)

		t, err := oauth2.ReuseTokenSource(z.token, src).Token()
		if err != nil {
			return err
		}

		z.token = t
	}

	return nil
}
