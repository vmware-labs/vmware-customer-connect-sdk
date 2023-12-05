// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache 2.0

package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDetailsSuccess(t *testing.T) {
	var dlgDetails DlgDetails
	dlgDetails, err = basicClient.GetDlgDetails("VMTOOLS1130", "1073")
	assert.Nil(t, err)
	assert.NotEmpty(t, dlgDetails.DownloadDetails, "Expected response to no be empty")
}

func TestGetDetailsInvalidProductId(t *testing.T) {
	var dlgDetails DlgDetails
	dlgDetails, err = basicClient.GetDlgDetails("VMTOOLS1130", "6666666")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrorDlgDetailsInputs)
	assert.Empty(t, dlgDetails, "Expected response to be empty")
}

func TestGetDetailsInvalidDownloadGroup(t *testing.T) {
	var dlgDetails DlgDetails
	dlgDetails, err = basicClient.GetDlgDetails("VMTOOLS666", "1073")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrorDlgDetailsInputs)
	assert.Empty(t, dlgDetails, "Expected response to be empty")
}

func TestFindDlgDetailsSuccess(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadDetails FoundDownload
	downloadDetails, err = authenticatedClient.FindDlgDetails("VMTOOLS1130", "1073", "VMware-Tools-darwin-*.tar.gz")
	assert.Nil(t, err)
	require.NotEmpty(t, downloadDetails.DownloadDetails)
	assert.NotEmpty(t, downloadDetails.DownloadDetails[0].FileName, "Expected response to not be empty")
}

func TestFindDlgDetailsGlobMultipleResults(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadDetails FoundDownload
	downloadDetails, err = authenticatedClient.FindDlgDetails("VMTOOLS1130", "1073", "*")
	assert.Nil(t, err)
	assert.Greater(t, len(downloadDetails.DownloadDetails), 1, "Expected response to be empty")
}

func TestFindDlgDetailsMultipleGlob(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadDetails FoundDownload
	downloadDetails, err = authenticatedClient.FindDlgDetails("VMTOOLS1130", "1073", "VMware-Tools-*-core-offline-depot-ESXi-all-*.zip")
	assert.Nil(t, err)
	require.NotEmpty(t, downloadDetails.DownloadDetails)
	assert.NotEmpty(t, downloadDetails.DownloadDetails[0].FileName, "Expected response to not be empty")
}

func TestFindDlgDetailsNoGlob(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadDetails FoundDownload
	downloadDetails, err = authenticatedClient.FindDlgDetails("VMTOOLS1130", "1073", "VMware-Tools-11.3.0-core-offline-depot-ESXi-all-18090558.zip")
	assert.Nil(t, err)
	require.NotEmpty(t, downloadDetails.DownloadDetails)
	assert.NotEmpty(t, downloadDetails.DownloadDetails[0].FileName, "Expected response to not be empty")
}

func TestFindDlgDetailsNoMatch(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadDetails FoundDownload
	downloadDetails, err = authenticatedClient.FindDlgDetails("VMTOOLS1130", "1073", "invalid*glob")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrorNoMatchingFiles)
	assert.Empty(t, downloadDetails.DownloadDetails, "Expected response to be empty")
}

func TestGetFileArray(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var fileArray []string
	fileArray, err = authenticatedClient.GetFileArray("vmware_horizon", "dem+standard", "2106", "PRODUCT_BINARY")
	assert.Nil(t, err)
	assert.NotEmpty(t, fileArray, "Expected response to no be empty")
}

func TestGetGetDlgProduct(t *testing.T) {
	var productID string
	var apiVersions APIVersions
	productID, apiVersions, err = basicClient.GetDlgProduct("vmware_tools", "vmtools", "11.1.1", "PRODUCT_BINARY")
	assert.Nil(t, err)
	assert.NotEmpty(t, apiVersions.Code, "Expected response to no be empty")
	assert.NotEmpty(t, productID, "Expected response to no be empty")
}

func TestGetGetDlgProductNsx(t *testing.T) {
	var productID string
	var apiVersions APIVersions
	productID, apiVersions, err = basicClient.GetDlgProduct("vmware_nsx", "nsx", "4.0*", "PRODUCT_BINARY")
	assert.Nil(t, err)
	assert.NotContains(t, apiVersions.Code, "-LE")
	assert.NotEmpty(t, productID, "Expected response to no be empty")
}

func TestGetGetDlgProductNsxLe(t *testing.T) {
	var productID string
	var apiVersions APIVersions
	productID, apiVersions, err = basicClient.GetDlgProduct("vmware_nsx", "nsx_le", "4.0.1.1 LE", "PRODUCT_BINARY")
	assert.Nil(t, err)
	assert.Contains(t, apiVersions.Code, "-LE")
	assert.NotEmpty(t, productID, "Expected response to no be empty")
}

func TestGetGetDlgProductNsxT(t *testing.T) {
	var productID string
	var apiVersions APIVersions
	productID, apiVersions, err = basicClient.GetDlgProduct("vmware_nsx_t_data_center", "nsx-t", "3.2*", "PRODUCT_BINARY")
	assert.Nil(t, err)
	assert.NotContains(t, apiVersions.Code, "-LE")
	assert.NotEmpty(t, productID, "Expected response to no be empty")
}

func TestGetGetDlgProductNsxTLe(t *testing.T) {
	var productID string
	var apiVersions APIVersions
	productID, apiVersions, err = basicClient.GetDlgProduct("vmware_nsx_t_data_center", "nsx-t_le", "3.2.1.2 LE", "PRODUCT_BINARY")
	assert.Nil(t, err)
	assert.Contains(t, apiVersions.Code, "-LE")
	assert.NotEmpty(t, productID, "Expected response to no be empty")
}