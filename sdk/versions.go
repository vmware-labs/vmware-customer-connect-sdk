package sdk

import (
	"errors"
	"sort"
	"strings"
)

type APIVersions struct {
	Code         string
	MajorVersion string
}

var ErrorNoMatchingVersions = errors.New("versions: invalid glob. no versions found")
var ErrorNoVersinGlob = errors.New("versions: invalid glob. glob my be provided")
var ErrorMultipleVersionGlob = errors.New("versions: invalid glob. a single version glob must be used")

func (c *Client) GetVersionMap(slug, subProductName string) (data map[string]APIVersions, err error) {
	data = make(map[string]APIVersions)

	var subProductDetails  SubProductDetails
	subProductDetails, err = c.GetSubProduct(slug, subProductName)
	if err != nil {
		return
	}

	// Loop through each major version and collect all versions
	for majorVersion, dlgList := range subProductDetails.DlgListByVersion {
		var dlgHeader DlgHeader
			dlgHeader, err = c.GetDlgHeader(dlgList.Code, dlgList.ProductID)
			if err != nil {
				return
			}

			for _, version := range dlgHeader.Versions {
				aPIVersions := APIVersions{
					Code:         version.ID,
					MajorVersion: majorVersion,
				}
				data[version.Name] = aPIVersions
			}
	}

	return
}

func (c *Client) FindVersion(slug, subProduct, version string) (data APIVersions, err error) {
	var versionMap map[string]APIVersions
	versionMap, err = c.GetVersionMap(slug, subProduct)
	if err != nil {
		return
	}

	var searchVersion string
	if strings.Contains(version, "*") {
		searchVersion, err = c.FindVersionFromGlob(slug, subProduct, version, versionMap)
		if err != nil {
			return
		}
	} else {
		searchVersion = version
	}
	
	if _, ok := versionMap[searchVersion]; !ok {
		err = ErrorInvalidVersion
		return
	}

	data = versionMap[searchVersion]
	return
}

func (c *Client) FindVersionFromGlob(slug, subProduct, versionGlob string, versionMap map[string]APIVersions) (version string, err error) {
	// Ensure only one glob is defined
	globCount := strings.Count(versionGlob, "*")
	if globCount == 0 {
		err = ErrorNoVersinGlob
		return
	} else if globCount > 1 {
		err = ErrorMultipleVersionGlob
		return
	}

	// Extract prefix by removing *
	versionPrefix := strings.Split(versionGlob, "*")[0]

	sortedKeys := sortVersionMapKeys(versionMap)

	// Check if only * is provided as strings.Split returns empty if separator is found
	if versionPrefix == "" {
		// return the first entry, which is the highest number.
		version = sortedKeys[0]
		return
	} else {
		// return the first entry matching the prefix
		for _, key := range sortedKeys {
			if strings.HasPrefix(key, versionPrefix) {
				version = key
				return
			}
		}
	}

	err = ErrorNoMatchingVersions
	return
}

func (c *Client) GetVersionSlice(slug, subProductName string) (data []string, err error) {
	var versionMap map[string]APIVersions
	versionMap, err = c.GetVersionMap(slug, subProductName)
	if err != nil {
		return
	}

	data = sortVersionMapKeys(versionMap)

	return
}

func sortVersionMapKeys(versionMap map[string]APIVersions) (keys []string) {
	// Extract all keys which are the version strings and reverse sort them
	// This means the versions will go from high to low
	keys = make([]string, len(versionMap))
	i := 0
	for key := range versionMap {
		keys[i] = key
		i++
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	return
}
