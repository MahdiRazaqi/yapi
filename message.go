package zoho

import (
	"fmt"
	"net/http"
)

// SendMail send email with zoho mail
func (z *Zoho) SendMail(from, to, subject, content string) error {
	params := parameters{
		"token":       z.token.AccessToken,
		"fromAddress": from,
		"toAddress":   to,
		"subject":     subject,
		"content":     content,
		"askReceipt":  "yes",
	}

	path := fmt.Sprintf("/accounts/%v/messages", z.accountID)

	_, err := request(path, http.MethodPost, params)
	return err
}
