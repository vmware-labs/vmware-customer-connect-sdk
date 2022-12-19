package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	majorVersionsURL = baseURL + "/channel/public/api/v1.0/products/getProductHeader"
)

type ProductVersions struct {
	MajorVersions []MajorVersions `json:"versions"`
	Resources     []interface{}   `json:"resources"`
}
type MajorVersions struct {
	ID string `json:"id"`
}

var ErrorInvalidMajorVersion = errors.New("major version not found")

// curl "https://customerconnect.vmware.com/channel/public/api/v1.0/products/getProductHeader?category=datacenter_cloud_infrastructure&product=vmware_vsphere_storage_appliance&version=5_5" |jq
func (c *Client) GetMajorVersionsSlice(slug string) (data []string, err error) {
	c.EnsureProductDetailMap()
	if err != nil {
		return
	}

	if _, ok := ProductDetailMap[slug]; !ok {
		err = ErrorInvalidSlug
		return
	}

	search_string := fmt.Sprintf("?category=%s&product=%s&version=%s",
		ProductDetailMap[slug].Category, slug, ProductDetailMap[slug].LatestMajorVersion)

	res, err := c.HttpClient.Get(majorVersionsURL + search_string)
	if err != nil {
		return
	}
	defer res.Body.Close()

	var productVersions ProductVersions
	err = json.NewDecoder(res.Body).Decode(&productVersions)

	if err == nil {
		for _, version := range productVersions.MajorVersions {
			data = append(data, version.ID)
		}
	}

	return
}
