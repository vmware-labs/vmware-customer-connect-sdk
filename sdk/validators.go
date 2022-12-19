package sdk

import (
	"errors"
	"net/http"
)

var ErrorInvalidSlug = errors.New("api: slug is not valid")
var ErrorInvalidCategory = errors.New("api: category is not valid")
var ErrorInvalidVersion = errors.New("api: version is not valid")
var ErrorServerError = errors.New("api: server down. 500 error received")

func (c *Client) validateSlugCategoryVersion(slug, category, majorVersion string) (err error) {
	c.EnsureProductDetailMap()
	if err != nil {
		return
	}

	if _, ok := ProductDetailMap[slug]; !ok {
		err = ErrorInvalidSlug
		return
	}

	if ProductDetailMap[slug].LatestMajorVersion != majorVersion {
		err = ErrorInvalidVersion
		return
	}
	return
	// TODO add check for non-latest version
	// Potential for circular dependency as this validator used by version get command
}

func (c *Client) validateResponseSlugCategoryVersion(slug, category, majorVersion string, res http.Response) (err error) {
	if res.StatusCode == 400 {
		err = c.validateSlugCategoryVersion(slug, category, majorVersion)
		return
	}
	return
}

func (c *Client) validateResponseGeneric(resCode int) (err error) {
	if resCode == 401 {
		err = ErrorNotAuthorized
		return
	} else if resCode == 500 {
		err = ErrorServerError
		return
	} else if resCode != 200 {
		err = ErrorNon200Response
		return
	}
	return
}
