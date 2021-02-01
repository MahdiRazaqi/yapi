package zoho

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

// Zoho model
type Zoho struct {
	Config    *oauth2.Config
	Token     *oauth2.Token
	AccountID string
}

// New zoho model
func New(clientID, clientSecret, redirectURL, authURL, tokenURL string, scopes []string) *Zoho {
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		RedirectURL:  redirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:   authURL,
			TokenURL:  tokenURL,
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	url := fmt.Sprintf("https://accounts.zoho.com/oauth/v2/auth?scope=%v&client_id=%v&response_type=code&access_type=offline&redirect_uri=%v", strings.Join(scopes, ","), clientID, redirectURL)
	fmt.Printf("Open this link for authentication:\n%v\n\n", url)

	return &Zoho{
		Config: conf,
	}
}

// AuthWithCode get token from code
func (z *Zoho) AuthWithCode(code string) error {
	t, err := z.Config.Exchange(context.Background(), code)
	if err != nil {
		return err
	}
	z.Token = t

	resp, err := request("/accounts", http.MethodGet, parameters{
		"token": z.Token.AccessToken,
	})
	if err != nil {
		return err
	}

	z.AccountID = resp.toAccount()[0].AccountID

	return nil
}

// RenewToken renew token
func (z *Zoho) RenewToken() error {
	src := z.Config.TokenSource(oauth2.NoContext, z.Token)

	t, err := oauth2.ReuseTokenSource(z.Token, src).Token()
	if err != nil {
		return err
	}

	z.Token = t

	return nil
}
