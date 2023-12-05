// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache 2.0

package sdk

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSubProductsSlice(t *testing.T) {
	var subProducts []SubProductDetails
	subProducts, err = basicClient.GetSubProductsSlice("vmware_horizon", "PRODUCT_BINARY")
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(subProducts), 20, "Expected response to contain at least 20 items")
}

func TestGetSubProductsSliceDrivers(t *testing.T) {
	var subProducts []SubProductDetails
	subProducts, err = basicClient.GetSubProductsSlice("vmware_vsphere", "DRIVERS_TOOLS")
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(subProducts), 20, "Expected response to contain at least 20 items")
}

func TestGetSubProductNsxLE(t *testing.T) {
	var subProduct SubProductDetails
	subProduct, err = basicClient.GetSubProduct("vmware_nsx_t_data_center", "nsx-t_le", "PRODUCT_BINARY")
	assert.Nil(t, err)
	assert.NotEmpty(t, subProduct.ProductName)
	assert.Greater(t, len(subProduct.DlgListByVersion), 0)
}

// func TestGetSubProduct(t *testing.T) {
// 	var subProduct SubProductDetails
// 	subProduct, err = basicClient.GetSubProduct("vmware_vsphere", "dem+standard", "PRODUCT_BINARY")
// 	assert.Nil(t, err)
// 	assert.NotEmpty(t, subProduct.ProductName)
// }

func TestGetSubProductDriver(t *testing.T) {
	var subProduct SubProductDetails
	subProduct, err = basicClient.GetSubProduct("vmware_horizon", "dem+standard", "PRODUCT_BINARY")
	assert.Nil(t, err)
	assert.NotEmpty(t, subProduct.ProductName)
}

func TestGetSubProductsSliceInvalidSlug(t *testing.T) {
	var subProducts []SubProductDetails
	subProducts, err = basicClient.GetSubProductsSlice("vsphere", "PRODUCT_BINARY")
	assert.ErrorIs(t, err, ErrorInvalidSlug)
	assert.Empty(t, subProducts, "Expected response to be empty")
}

func TestGetSubProductsMap(t *testing.T) {
	var subProducts map[string]SubProductDetails
	subProducts, err = basicClient.GetSubProductsMap("vmware_vsphere", "PRODUCT_BINARY")
	assert.Nil(t, err)
	assert.Contains(t, subProducts, "vmtools")
}

func TestGetSubProductsMapHorizon(t *testing.T) {
	var subProducts map[string]SubProductDetails
	subProducts, err = basicClient.GetSubProductsMap("vmware_horizon_clients", "PRODUCT_BINARY")
	assert.Nil(t, err)
	assert.Contains(t, subProducts, "cart+win")
	assert.Contains(t, subProducts, "cart+andrd_x8632")
	assert.Contains(t, subProducts, "cart+lin64")
}

func TestGetSubProductsMapNsxLe(t *testing.T) {
	var subProducts map[string]SubProductDetails
	subProducts, err = basicClient.GetSubProductsMap("vmware_nsx", "PRODUCT_BINARY")
	assert.Nil(t, err)
	assert.Contains(t, subProducts, "nsx")
	assert.Contains(t, subProducts, "nsx_le")
}

func TestGetSubProductsMapNsxTLe(t *testing.T) {
	var subProducts map[string]SubProductDetails
	subProducts, err = basicClient.GetSubProductsMap("vmware_nsx_t_data_center", "PRODUCT_BINARY")
	assert.Nil(t, err)
	assert.Contains(t, subProducts, "nsx-t")
	assert.Contains(t, subProducts, "nsx-t_le")
}

func TestGetSubProductsMapInvalidSlug(t *testing.T) {
	var subProductMap map[string]SubProductDetails
	subProductMap, err = basicClient.GetSubProductsMap("vsphere", "PRODUCT_BINARY")
	assert.ErrorIs(t, err, ErrorInvalidSlug)
	assert.Empty(t, subProductMap, "Expected response to be empty")
}

func TestGetSubProductsDetails(t *testing.T) {
	var subProductDetails DlgList
	subProductDetails, err = basicClient.GetSubProductDetails("vmware_vsphere", "vmtools", "6_7", "PRODUCT_BINARY")
	assert.Nil(t, err)
	assert.NotEmpty(t, subProductDetails.Code, "Expected response to not be empty")
}

func TestGetSubProductsDetailsInvalidSubProduct(t *testing.T) {
	var subProductDetails DlgList
	subProductDetails, err = basicClient.GetSubProductDetails("vmware_vsphere", "tools", "6_7", "PRODUCT_BINARY")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrorInvalidSubProduct)
	assert.Empty(t, subProductDetails.Code, "Expected response to be empty")
}

func TestGetSubProductsDetailsInvalidMajorVersion(t *testing.T) {
	var subProductDetails DlgList
	subProductDetails, err = basicClient.GetSubProductDetails("vmware_vsphere", "vmtools", "5_5", "PRODUCT_BINARY")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrorInvalidSubProductMajorVersion)
	assert.Empty(t, subProductDetails.Code, "Expected response to be empty")
}

func TestModifyHorizonClientCode(t *testing.T) {
	productCode := "cart24fq4_lin_2309.1_tarball"
	productCode = modifyHorizonClientCode(productCode)
	assert.Equal(t, "cart+tarball", productCode)

	productCode = "one_2"
	productCode = modifyHorizonClientCode(productCode)
	assert.Equal(t, "one", productCode)
}

func TestGetProductName(t *testing.T) {
	reEndVersion := regexp.MustCompile(`[0-9]+.*`)
	productName := "VMware vSphere Hypervisor (ESXi) 8.0U2"
	productName = getProductName(productName, "vmware_vsphere", "", reEndVersion)
	assert.Equal(t, "VMware vSphere Hypervisor (ESXi)", productName)

	// Ensure drivers are unmodified
	productName = "VMware ESXi 8.0 native ixgben ENS 1.18.2.0 NIC Driver for Intel Ethernet Controllers 82599, x520, x540, x550, and x552 family"
	productName = getProductName(productName, "vmware_vsphere", "Driver CDs", reEndVersion)
	assert.Equal(t, productName, "Driver - native ixgben ENS")
}

func TestGetProductCode(t *testing.T) {
	reEndVersion := regexp.MustCompile(`[0-9]+.*`)
	productCode := "ESXI80U2"
	productCode = getProductCode(productCode, "vmware_vsphere", "", reEndVersion)
	assert.Equal(t, "esxi", productCode)

	// Ensure drivers are unmodified
	productCode = "DT-ESXI80-INTEL-I40EN-2650-1OEM"
	productCode = getProductCode(productCode, "vmware_vsphere", "Driver CDs", reEndVersion)
	assert.Equal(t, "dt-esxi80-intel-i40en-2650-1oem", productCode)
}