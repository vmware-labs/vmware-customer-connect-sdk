package sdk

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type AccountInfo struct {
	UserType    string   `json:"userType"`
	AccountList []AccntList `json:"accntList"`
}

type AccntList struct {
	EaNumber  string `json:"eaNumber"`
	EaName    string `json:"eaName"`
	IsDefault string `json:"isDefault"`
}

type CurrentUser struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

const (
	accountInfoURL = baseURL + "/channel/api/v1.0/ems/accountinfo"
	currentUserURL = baseURL + "/vmwauth/loggedinuser"
)

var ErrorNotAuthorized = errors.New("account: you are not authenticated")
var ErrorNon200Response = errors.New("account: server did not respond with 200 ok")

func (c *Client) AccountInfo() (data AccountInfo, err error) {
	payload := `{"rowLimit": 1000}`
	var res *http.Response
	res, err = c.HttpClient.Post(accountInfoURL, "application/json", strings.NewReader(payload))
	if err != nil {
		return
	}
	defer res.Body.Close()

	if err = c.validateResponseGeneric(res.StatusCode); err != nil {
		return
	}

	err = json.NewDecoder(res.Body).Decode(&data)

	return
}

func (c *Client) CheckLoggedIn() (err error) {
	_, err = c.AccountInfo()
	return
}

func (c *Client) CurrentUser() (data CurrentUser, err error) {
	if err = c.CheckLoggedIn(); err != nil {
		return
	}

	var res *http.Response
	res, err = c.HttpClient.Get(currentUserURL)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if err = c.validateResponseGeneric(res.StatusCode); err != nil {
		return
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	return
}
