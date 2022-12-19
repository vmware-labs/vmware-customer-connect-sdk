package sdk

import (
	// "fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSubProductsSlice(t *testing.T) {
	var subProducts []SubProductDetails
	subProducts, err = basicClient.GetSubProductsSlice("vmware_horizon")
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(subProducts), 20, "Expected response to contain at least 20 items")
}

func TestGetSubProduct(t *testing.T) {
	var subProduct SubProductDetails
	subProduct, err = basicClient.GetSubProduct("vmware_horizon", "dem+standard")
	assert.Nil(t, err)
	assert.NotEmpty(t, subProduct.ProductName)
}

// func TestGetSubProductsSliceAll(t *testing.T) {
// 	var products map[string]ProductDetails
// 	products, err = basicClient.GetProductsMap()

// 	var allSubProducts []string
// 	for product, _ := range products {
// 		fmt.Println(product)
// 		var subProductsStr []string
// 		var subProducts []SubProduct
// 		subProducts, err = basicClient.GetSubProductsSlice(product)
// 		for _, subProduct := range subProducts {
// 			combined_code := fmt.Sprintf("%s>>%s", product, subProduct.ProductCode)
// 			subProductsStr = append(subProductsStr, subProduct.ProductCode)
// 			allSubProducts = append(allSubProducts, combined_code)

// 		}
// 		fmt.Println(subProductsStr)
// 	}
// 	fmt.Println(allSubProducts)
// 	// t.Log(allSubProducts)
// 	assert.NotNil(t, err)
	
// }

func TestGetSubProductsSliceInvalidSlug(t *testing.T) {
	var subProducts []SubProductDetails
	subProducts, err = basicClient.GetSubProductsSlice("vsphere")
	assert.ErrorIs(t, err, ErrorInvalidSlug)
	assert.Empty(t, subProducts, "Expected response to be empty")
}

func TestGetSubProductsMap(t *testing.T) {
	var subProducts map[string]SubProductDetails
	subProducts, err = basicClient.GetSubProductsMap("vmware_vsphere")
	assert.Nil(t, err)
	assert.Contains(t, subProducts, "vmtools")
}

func TestGetSubProductsMapInvalidSlug(t *testing.T) {
	var subProductMap map[string]SubProductDetails
	subProductMap, err = basicClient.GetSubProductsMap("vsphere")
	assert.ErrorIs(t, err, ErrorInvalidSlug)
	assert.Empty(t, subProductMap, "Expected response to be empty")
}

func TestGetSubProductsDetails(t *testing.T) {
	var subProductDetails DlgList
	subProductDetails, err = basicClient.GetSubProductDetails("vmware_vsphere", "vmtools", "6_7")
	assert.Nil(t, err)
	assert.NotEmpty(t, subProductDetails.Code, "Expected response to not be empty")
}

func TestGetSubProductsDetailsInvalidSubProduct(t *testing.T) {
	var subProductDetails DlgList
	subProductDetails, err = basicClient.GetSubProductDetails("vmware_vsphere", "tools", "6_7")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrorInvalidSubProduct)
	assert.Empty(t, subProductDetails.Code, "Expected response to be empty")
}

func TestGetSubProductsDetailsInvalidMajorVersion(t *testing.T) {
	var subProductDetails DlgList
	subProductDetails, err = basicClient.GetSubProductDetails("vmware_vsphere", "vmtools", "5_5")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrorInvalidSubProductMajorVersion)
	assert.Empty(t, subProductDetails.Code, "Expected response to be empty")
}
