package sdk

import (
	"errors"
	// "fmt"
	// "reflect"
	"regexp"
	"sort"
	"strings"
)

type SubProductDetails struct {
	ProductName      string
	ProductCode      string
	DlgListByVersion map[string]DlgList
}

type SubProductSliceElement struct {
	Name        string
	Description string
}

var ErrorInvalidSubProduct = errors.New("subproduct: invalid subproduct requested")
var ErrorInvalidSubProductMajorVersion = errors.New("subproduct: invalid major version requested")

func (c *Client) GetSubProductsMap(slug string) (data map[string]SubProductDetails, err error) {
	c.EnsureProductDetailMap()
	if err != nil {
		return
	}

	if _, ok := ProductDetailMap[slug]; !ok {
		err = ErrorInvalidSlug
		return
	}

	subProductMap := make(map[string]SubProductDetails)

	var majorVersions []string
	majorVersions, err = c.GetMajorVersionsSlice(slug)
	if err != nil {
		return
	}

	// Iterate major product versions and extract all unique products
	// All version information is stripped
	for _, majorVersion := range majorVersions {
		var dlgEditionsList []DlgEditionsLists
		dlgEditionsList, err = c.GetDlgEditionsList(slug, majorVersion)
		// Invalid version errors need to be ignored, as they come from deprecated products
		if err == ErrorInvalidVersion {
			err = nil
			continue
		} else if err != nil {
			return
		}

		for _, dlgEdition := range dlgEditionsList {
			for _, dlgList := range dlgEdition.DlgList {
				productCode := strings.ToLower(dlgList.Code)
				productName := dlgList.Name
				// Regex captures numbers and all text after
				reEndVersion := regexp.MustCompile(`[0-9]+.*`)
				// Regex detects numbers surrounded by - or _
				reMidVersion := regexp.MustCompile(`(-|_)([0-9.]+)(-|_)`)

				// Horizon clients don't follow a common pattern for API naming. This block aligns the pattern
				if strings.HasPrefix(productCode, "cart") {
					productCode = strings.Replace(productCode, "-", "_", 1)
					
					// Remove version numbers at the start of the string only
					reHorizon := regexp.MustCompile(`([0-9.].*?)_`)
					found := reHorizon.FindString(productCode)
					if found != "" {
						productCode = strings.Replace(productCode, found, "+", 1)
					}
					// Handle tarball not following pattern. Replace cart+lin_+tarball to cart+tarball
					if strings.HasSuffix(productCode, "tarball") {
						// productCode = strings.Replace(productCode, "lin_+", "", 1)
						reHorizonTar := regexp.MustCompile(`lin_([0-9]+.*?)_`)
						productCode = reHorizonTar.ReplaceAllString(productCode, "")
					} else {
						// Remove version numbers at the end
						reHorizonVersion := regexp.MustCompile(`_([0-9.].*)`)
						productCode = reHorizonVersion.ReplaceAllString(productCode, "")
					}

				} else {
					// Check if product code has text after the version section
					if ok := reMidVersion.MatchString(productCode); ok{
						// replace version with + to allow for string to be split when searching
						productCode = reMidVersion.ReplaceAllString(productCode, "+")
						// remove versions prepended versions
						reFpStrip := regexp.MustCompile(`(\+fp[0-9])|(\+hf[0-9])`)
						productCode = reFpStrip.ReplaceAllString(productCode, "")
					} else {
						// when product ends with a version, remove all text after the first number
						productCode = reEndVersion.ReplaceAllString(productCode, "")
						productCode = strings.TrimSuffix(strings.TrimSuffix(productCode, "_"), "-")
					}

				}

				// Special case for Horizon due to inconsistent naming
				if slug == "vmware_horizon" {
					reNumbers := regexp.MustCompile(`[0-9.,]+`)
					reSpace := regexp.MustCompile(`\s+`)
					productName = strings.TrimSpace(reNumbers.ReplaceAllString(productName, ""))
					productName = reSpace.ReplaceAllString(productName, " ")
				} else {
					productName = strings.TrimSpace(reEndVersion.ReplaceAllString(productName, ""))
				}

				// Initalize the struct for a product code for the first time
				if _, ok := subProductMap[productCode]; !ok {
					subProductMap[productCode] = SubProductDetails{
						ProductName:      productName,
						ProductCode:      productCode,
						DlgListByVersion: make(map[string]DlgList),
					}
				}

				subProductMap[productCode].DlgListByVersion[majorVersion] = dlgList
			}
		}
	}

	data = subProductMap
	return
}


func (c *Client) GetSubProductsSlice(slug string) (data []SubProductDetails, err error) {
	subProductMap, err := c.GetSubProductsMap(slug)
	if err != nil {
		return
	}

	// Sort keys to output sorted slice
	keys := make([]string, len(subProductMap))
	i := 0
	for key := range subProductMap {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	// Append to array using sorted keys to fetch from map
	for _, key := range keys {
		data = append(data, subProductMap[key])
	}

	return
}

func (c *Client) GetSubProduct(slug, subProduct string) (data SubProductDetails, err error) {
	var subProductMap map[string]SubProductDetails
	subProductMap, err = c.GetSubProductsMap(slug)
	if err != nil {
		return
	}

	if foundSubProduct, ok := subProductMap[subProduct]; !ok {
		err = ErrorInvalidSubProduct
	} else {
		data = foundSubProduct
	}

	return
}

func (c *Client) GetSubProductDetails(slug, subProduct, majorVersion string) (data DlgList, err error) {
	var subProducts map[string]SubProductDetails
	subProducts, err = c.GetSubProductsMap(slug)
	if err != nil {
		return
	}

	if subProduct, ok := subProducts[subProduct]; ok {
		if dlgList, ok := subProduct.DlgListByVersion[majorVersion]; ok {
			data = dlgList
		} else {
			err = ErrorInvalidSubProductMajorVersion
		}

	} else {
		err = ErrorInvalidSubProduct
	}

	return
}
