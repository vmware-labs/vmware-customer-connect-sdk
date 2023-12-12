// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache 2.0

package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchDownloadLinkVersionGlob(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadPayload []DownloadPayload
	downloadPayload, err = authenticatedClient.GenerateDownloadPayload("vmware_tools", "vmtools", "11.*", "VMware-Tools-darwin-*.tar.gz", "PRODUCT_BINARY", true)
	require.Nil(t, err)
	require.NotEmpty(t, downloadPayload)
	assert.NotEmpty(t, downloadPayload[0].ProductId, "Expected response not to be empty")

	t.Logf(fmt.Sprintf("download_payload: %+v\n", downloadPayload))

	var authorizedDownload AuthorizedDownload
	authorizedDownload, _ = authenticatedClient.FetchDownloadLink(downloadPayload[0])
	assert.Nil(t, err)
	assert.NotEmpty(t, authorizedDownload.DownloadURL, "Expected response not to be empty")

	t.Logf(fmt.Sprintf("download_details: %+v\n", authorizedDownload))
}

func TestFetchDownloadLinkVersionDrivers(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadPayload []DownloadPayload
	downloadPayload, err = authenticatedClient.GenerateDownloadPayload("vmware_vsphere", "vs-mgmt-sdk80u2", "8.0U2", "VMware-vSphere-SDK-8.0.2-22394481.zip", "DRIVERS_TOOLS", true)
	require.Nil(t, err)
	require.NotEmpty(t, downloadPayload)
	assert.NotEmpty(t, downloadPayload[0].ProductId, "Expected response not to be empty")

	t.Logf(fmt.Sprintf("download_payload: %+v\n", downloadPayload))

	var authorizedDownload AuthorizedDownload
	authorizedDownload, _ = authenticatedClient.FetchDownloadLink(downloadPayload[0])
	assert.Nil(t, err)
	assert.NotEmpty(t, authorizedDownload.DownloadURL, "Expected response not to be empty")

	t.Logf(fmt.Sprintf("download_details: %+v\n", authorizedDownload))
}

func TestFetchDownloadLinkVersionCustomIso(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadPayload []DownloadPayload
	downloadPayload, err = authenticatedClient.GenerateDownloadPayload("vmware_vsphere", "oem-esxi80u2-hitachi", "ESXi 8.0.2", "VMware-ESXi-8.0-update2-22380479-hitachi-1301.iso", "CUSTOM_ISO", true)
	require.Nil(t, err)
	require.NotEmpty(t, downloadPayload)
	assert.NotEmpty(t, downloadPayload[0].ProductId, "Expected response not to be empty")

	t.Logf(fmt.Sprintf("download_payload: %+v\n", downloadPayload))

	var authorizedDownload AuthorizedDownload
	authorizedDownload, _ = authenticatedClient.FetchDownloadLink(downloadPayload[0])
	assert.Nil(t, err)
	assert.NotEmpty(t, authorizedDownload.DownloadURL, "Expected response not to be empty")

	t.Logf(fmt.Sprintf("download_details: %+v\n", authorizedDownload))
}

func TestFetchDownloadLinkVersionAddons(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadPayload []DownloadPayload
	downloadPayload, err = authenticatedClient.GenerateDownloadPayload("vmware_vsphere", "addon_esxi80u2_hitachi", "ESXi 8.0U2", "VMware-ESXi-8.0u2-addon-22380479-hitachi-1301.zip", "ADDONS", true)
	require.Nil(t, err)
	require.NotEmpty(t, downloadPayload)
	assert.NotEmpty(t, downloadPayload[0].ProductId, "Expected response not to be empty")

	t.Logf(fmt.Sprintf("download_payload: %+v\n", downloadPayload))

	var authorizedDownload AuthorizedDownload
	authorizedDownload, _ = authenticatedClient.FetchDownloadLink(downloadPayload[0])
	assert.Nil(t, err)
	assert.NotEmpty(t, authorizedDownload.DownloadURL, "Expected response not to be empty")

	t.Logf(fmt.Sprintf("download_details: %+v\n", authorizedDownload))
}

func TestFetchDownloadLinkVersion(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadPayload []DownloadPayload
	downloadPayload, err = authenticatedClient.GenerateDownloadPayload("vmware_tools", "vmtools", "11.1.1", "VMware-Tools-darwin-*.tar.gz", "PRODUCT_BINARY", true)
	assert.Nil(t, err)
	require.NotEmpty(t, downloadPayload)
	assert.NotEmpty(t, downloadPayload[0].ProductId, "Expected response not to be empty")

	t.Logf(fmt.Sprintf("download_payload: %+v\n", downloadPayload))

	var authorizedDownload AuthorizedDownload
	authorizedDownload, _ = authenticatedClient.FetchDownloadLink(downloadPayload[0])
	assert.Nil(t, err)
	assert.NotEmpty(t, authorizedDownload.DownloadURL, "Expected response not to be empty")

	t.Logf(fmt.Sprintf("download_details: %+v\n", authorizedDownload))
}

func TestFetchDownloadLinkInvalidVersion(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadPayload []DownloadPayload
	downloadPayload, err = authenticatedClient.GenerateDownloadPayload("vmware_tools", "vmtools", "666", "VMware-Tools-darwin-*.tar.gz", "PRODUCT_BINARY", true)
	assert.ErrorIs(t, err, ErrorInvalidVersion)
	assert.Empty(t, downloadPayload, "Expected response to be empty")
}

func TestFetchDownloadLinkNeedEula(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadPayload []DownloadPayload
	downloadPayload, err = authenticatedClient.GenerateDownloadPayload("vmware_tools", "vmtools", "11.1.0", "VMware-Tools-darwin-*.tar.gz", "PRODUCT_BINARY", false)
	assert.ErrorIs(t, err, ErrorEulaUnaccepted)
	assert.Empty(t, downloadPayload, "Expected response to be empty")
}

func TestFetchDownloadLinkNotEntitled(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadPayload []DownloadPayload
	downloadPayload, err = authenticatedClient.GenerateDownloadPayload("vmware_nsx_t_data_center", "nsx-t", "3.2.3.1", "nsx-unified-appliance-secondary-*.qcow2", "PRODUCT_BINARY", true)
	assert.ErrorIs(t, err, ErrorNotEntitled)
	assert.Empty(t, downloadPayload, "Expected response to be empty")
}

func TestGenerateDownloadInvalidVersionGlob(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadPayload []DownloadPayload
	downloadPayload, err = authenticatedClient.GenerateDownloadPayload("vmware_tools", "vmtools", "666.*", "VMware-Tools-darwin-*.tar.gz", "PRODUCT_BINARY", true)
	assert.ErrorIs(t, err, ErrorNoMatchingVersions)
	assert.Empty(t, downloadPayload, "Expected response to be empty")
}

func TestGenerateDownloadDoubleVersion(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadPayload []DownloadPayload
	downloadPayload, err = authenticatedClient.GenerateDownloadPayload("vmware_tools", "vmtools", "*.*", "VMware-Tools-darwin-*.tar.gz", "PRODUCT_BINARY", true)
	assert.ErrorIs(t, err, ErrorMultipleVersionGlob)
	assert.Empty(t, downloadPayload, "Expected response to be empty")
}
