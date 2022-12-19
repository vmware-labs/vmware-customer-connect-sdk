package sdk

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	eulaURL = baseURL + "/channel/api/v1.0/dlg/eula/accept"
)

var ErrorEulaInputs = errors.New("eula: downloadGroup or productId invalid")

func (c *Client) FetchEulaUrl(downloadGroup, productId string) (url string, err error) {
	if err = c.CheckLoggedIn(); err != nil {
		return
	}

	var dlgDetails DlgDetails
	dlgDetails, err = c.GetDlgDetails(downloadGroup, productId)
	if err != nil {
		return
	}

	url = dlgDetails.EulaResponse.EulaURL

	return
}

func (c *Client) AcceptEula(downloadGroup, productId string) (err error) {
	if err = c.CheckLoggedIn(); err != nil {
		return
	}

	search_string := fmt.Sprintf("?downloadGroup=%s&productId=%s", downloadGroup, productId)
	var res *http.Response
	res, err = c.HttpClient.Get(eulaURL + search_string)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode == 400 {
		err = ErrorEulaInputs
		return
	} else if res.StatusCode != 200 {
		err = ErrorNon200Response
	}

	return
}
