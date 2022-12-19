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
	fileArray, err = authenticatedClient.GetFileArray("vmware_horizon", "dem+standard", "2106")
	assert.Nil(t, err)
	assert.NotEmpty(t, fileArray, "Expected response to no be empty")
}

func TestGetGetDlgProduct(t *testing.T) {
	var downloadGroup, productID string
	downloadGroup, productID, err = basicClient.GetDlgProduct("vmware_tools", "vmtools", "11.1.1")
	assert.Nil(t, err)
	assert.NotEmpty(t, downloadGroup, "Expected response to no be empty")
	assert.NotEmpty(t, productID, "Expected response to no be empty")
}