package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type DlgHeader struct {
	Versions  []Versions    `json:"versions"`
	Product   Product       `json:"product"`
	Dlg       Dlg           `json:"dlg"`
	Resources []interface{} `json:"resources"`
}
type Versions struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	IsSelected bool   `json:"isSelected"`
}
type Product struct {
	ID               string `json:"id"`
	ReleasePackageID string `json:"releasePackageId"`
	Categorymap      string `json:"categorymap"`
	Productmap       string `json:"productmap"`
	Versionmap       string `json:"versionmap"`
	Name             string `json:"name"`
	Version          string `json:"version"`
}
type Dlg struct {
	Name           string `json:"name"`
	ReleaseDate    string `json:"releaseDate"`
	Type           string `json:"type"`
	Code           string `json:"code"`
	Documentation  string `json:"documentation"`
	InternalType   string `json:"internalType"`
	IsFreeProduct  bool   `json:"isFreeProduct"`
	IsThirdParty   bool   `json:"isThirdParty"`
	IsMassMarket   bool   `json:"isMassMarket"`
	TagID          int    `json:"tagId"`
	Notes          string `json:"notes"`
	Description    string `json:"description"`
	CompatibleWith string `json:"compatibleWith"`
}

const (
	dlgHeaderURL = baseURL + "/channel/public/api/v1.0/products/getDLGHeader"
)

var ErrorDlgHeader = errors.New("dlgHeader: downloadGroup or productId invalid")

// curl "https://my.vmware.com/channel/public/api/v1.0/products/getDLGHeader?downloadGroup=VMTOOLS1130&productId=1073" |jq
func (c *Client) GetDlgHeader(downloadGroup, productId string) (data DlgHeader, err error) {
	search_string := fmt.Sprintf("?downloadGroup=%s&productId=%s", downloadGroup, productId)
	var res *http.Response
	res, err = c.HttpClient.Get(dlgHeaderURL + search_string)
	if err != nil {return}
	defer res.Body.Close()

	if res.StatusCode == 400 {
		err = ErrorDlgHeader
		return
	}

	err = json.NewDecoder(res.Body).Decode(&data)

	return
}
