package sdk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetProducts(t *testing.T) {
	var products []MajorProducts
	products, err = basicClient.GetProductsSlice()
	assert.Nil(t, err)
	assert.Greater(t, len(products), 80, "Expected response to contain at least 80 items")
}

func TestGetProductMap(t *testing.T) {
	var products map[string]ProductDetails
	products, err = basicClient.GetProductsMap()
	assert.Nil(t, err)
	assert.Contains(t, products, "vmware_tools")
}
