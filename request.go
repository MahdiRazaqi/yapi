package zoho

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type parameters map[string]interface{}

type response struct {
	Status struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
	} `json:"status"`
	Data interface{} `json:"data"`
}

const endpoint = "https://mail.zoho.com/api"

func request(path, method string, params parameters) (*response, error) {
	var body io.Reader
	token, ok := params["token"].(string)
	if !ok {
		return nil, errors.New("token not found in params")
	}

	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	token = "Bearer " + token
	delete(params, "token")

	if method == http.MethodGet {
		q := u.Query()
		for k, v := range params {
			q.Set(k, fmt.Sprintf("%v", v))
		}

		u.RawQuery = q.Encode()
	} else {
		data, _ := json.Marshal(params)
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, endpoint+u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	r := &response{}
	if err := json.Unmarshal(data, r); err != nil {
		fmt.Println(string(data))
		return nil, err
	}

	if resp.StatusCode > 300 {
		return r, errors.New(string(data))
	}

	return r, nil
}
