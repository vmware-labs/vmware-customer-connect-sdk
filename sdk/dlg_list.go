package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// dlg list
type DlgEditions struct {
	DlgEditionsLists []DlgEditionsLists `json:"dlgEditionsLists"`
}
type DlgList struct {
	Name             string `json:"name"`
	Code             string `json:"code"`
	ProductID        string `json:"productId"`
	ReleaseDate      string `json:"releaseDate"`
	ReleasePackageID string `json:"releasePackageId"`
}
type DlgEditionsLists struct {
	Name    string    `json:"name"`
	DlgList []DlgList `json:"dlgList"`
	OrderID int       `json:"orderId"`
}

const (
	dlgListURL = baseURL + "/channel/public/api/v1.0/products/getRelatedDLGList"
)

// curl "https://my.vmware.com/channel/public/api/v1.0/products/getRelatedDLGList?category=datacenter_cloud_infrastructure&product=vmware_vsan&version=7_0&dlgType=PRODUCT_BINARY" |jq
func (c *Client) GetDlgEditionsList(slug, majorVersion string) (data []DlgEditionsLists, err error) {
	var category string
	category, err = c.GetCategory(slug)
	if err != nil {return}

	search_string := fmt.Sprintf("?category=%s&product=%s&version=%s&dlgType=PRODUCT_BINARY", category, slug, majorVersion)
	var res *http.Response
	res, err = c.HttpClient.Get(dlgListURL + search_string)
	if err != nil {return}
	defer res.Body.Close()

	err = c.validateResponseSlugCategoryVersion(slug, category, majorVersion, *res)
	if err != nil {return}

	var dlgEditions DlgEditions
	err = json.NewDecoder(res.Body).Decode(&dlgEditions)
	data = dlgEditions.DlgEditionsLists

	return
}
