package zoho

import (
	"encoding/json"
	"net/http"
)

type account struct {
	AccountID   string `json:"accountId"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DisplayName string `json:"displayName"`
	AccountName string `json:"accountName"`
	PhoneNumer  string `json:"phoneNumer"`
}

func (r *response) toAccount() (result []account) {
	data, err := json.Marshal(r.Data)
	if err != nil {
		return
	}

	if err = json.Unmarshal(data, &result); err != nil {
		return
	}

	return result
}

func (z *Zoho) setAccountID() error {
	resp, err := request("/accounts", http.MethodGet, parameters{
		"token": z.token.AccessToken,
	})
	if err != nil {
		return err
	}

	z.accountID = resp.toAccount()[0].AccountID

	return nil
}
