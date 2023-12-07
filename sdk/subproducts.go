// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache 2.0

package sdk

import (
	"errors"
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

func (c *Client) GetSubProductsMap(slug, dlgType string) (subProductMap map[string]SubProductDetails, err error) {
	c.EnsureProductDetailMap()
	if err != nil {
		return
	}

	if _, ok := ProductDetailMap[slug]; !ok {
		err = ErrorInvalidSlug
		return
	}

	subProductMap = make(map[string]SubProductDetails)

	var majorVersions []string
	majorVersions, err = c.GetMajorVersionsSlice(slug)
	if err != nil {
		return
	}

	// Iterate major product versions and extract all unique products
	// All version information is stripped
	for _, majorVersion := range majorVersions {
		var dlgEditionsList []DlgEditionsLists
		dlgEditionsList, err = c.GetDlgEditionsList(slug, majorVersion, dlgType)
		// Invalid version errors need to be ignored, as they come from deprecated products
		if err == ErrorInvalidVersion {
			err = nil
			continue
		} else if err != nil {
			return
		}

		for _, dlgEdition := range dlgEditionsList {
			for _, dlgList := range dlgEdition.DlgList {
				// Regex captures numbers and all text after
				reEndVersion := regexp.MustCompile(`[0-9]+.*`)
				
				productName := getProductName(dlgList.Name, slug, dlgType, reEndVersion)
				productCode := getProductCode(strings.ToLower(dlgList.Code), slug, dlgType, reEndVersion)

				// Initalize the struct for a product code for the first time
				if _, ok := subProductMap[productCode]; !ok {
					subProductMap[productCode] = SubProductDetails{
						ProductName:      productName,
						ProductCode:      productCode,
						DlgListByVersion: make(map[string]DlgList),
					}
				}

				subProductMap[productCode].DlgListByVersion[majorVersion] = dlgList

				if productCode == "nsx" || productCode == "nsx-t" {
					duplicateNsxToNsxLe(subProductMap, productCode, productName, majorVersion, dlgList)
				}
			}
		}
	}
	return
}

func getProductCode(productCode, slug, dlgType string, reEndVersion *regexp.Regexp) (string) {
	productCode = strings.ToLower(productCode)
	if dlgType != "PRODUCT_BINARY" {
		return productCode
	}
	reMidVersion := regexp.MustCompile(`(-|_)([0-9.]+)(-|_)`)
	if strings.HasPrefix(productCode, "cart") {
		return modifyHorizonClientCode(productCode)
		} 
	// Horizon clients don't follow a common pattern for API naming. This block aligns the pattern
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
	return productCode
}

func getProductName(productName, slug, dlgType string, reEndVersion *regexp.Regexp) (string) {
	if dlgType != "PRODUCT_BINARY" {
		return productName
	}
	// Special case for Horizon due to inconsistent naming
	if slug == "vmware_horizon" {
		reNumbers := regexp.MustCompile(`[0-9.,]+`)
		reSpace := regexp.MustCompile(`\s+`)
		productName := strings.TrimSpace(reNumbers.ReplaceAllString(productName, ""))
		return reSpace.ReplaceAllString(productName, " ")
	} else {
		return strings.TrimSpace(reEndVersion.ReplaceAllString(productName, ""))
	}
}

func modifyHorizonClientCode(productCode string) (string) {
	productCode = strings.Replace(productCode, "-", "_", 1)
					
	// Remove version numbers at the start of the string only
	reHorizon := regexp.MustCompile(`([0-9.].*?)_`)
	found := reHorizon.FindString(productCode)
	if found != "" {
		productCode = strings.Replace(productCode, found, "+", 1)
	}
	// Handle tarball not following pattern.
	if strings.HasSuffix(productCode, "tarball") {
		// productCode = strings.Replace(productCode, "lin_+", "", 1)
		reHorizonTar := regexp.MustCompile(`lin_([0-9]+.*?)_`)
		productCode = reHorizonTar.ReplaceAllString(productCode, "")
	} else {
		// Remove version numbers at the end
		reHorizonVersion := regexp.MustCompile(`_([0-9.].*)`)
		productCode = reHorizonVersion.ReplaceAllString(productCode, "")
	}
	return productCode
}

// Duplicate NSX LE to a separate subproduct
func duplicateNsxToNsxLe(subProductMap map[string]SubProductDetails, productCode, productName, majorVersion string, dlgList DlgList) {
	if _, ok := subProductMap[productCode + "_le"]; !ok {
		subProductMap[productCode + "_le"] = SubProductDetails{
			ProductName:      productName + " Limited Edition",
			ProductCode:      productCode + "_le",
			DlgListByVersion: make(map[string]DlgList),
		}
	}
	dlgList.Name = dlgList.Name + " Limited Edition"
	dlgList.Code = dlgList.Code + "-LE"
	subProductMap[productCode + "_le"].DlgListByVersion[majorVersion] = dlgList
}


func (c *Client) GetSubProductsSlice(slug, dlgType string) (data []SubProductDetails, err error) {
	subProductMap, err := c.GetSubProductsMap(slug, dlgType)
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

func (c *Client) GetSubProduct(slug, subProduct, dlgType string) (data SubProductDetails, err error) {
	var subProductMap map[string]SubProductDetails
	subProductMap, err = c.GetSubProductsMap(slug, dlgType)
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

func (c *Client) GetSubProductDetails(slug, subProduct, majorVersion, dlgType string) (data DlgList, err error) {
	var subProducts map[string]SubProductDetails
	subProducts, err = c.GetSubProductsMap(slug, dlgType)
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
