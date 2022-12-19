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
	downloadPayload, err = authenticatedClient.GenerateDownloadPayload("vmware_tools", "vmtools", "11.*", "VMware-Tools-darwin-*.tar.gz", true)
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
	downloadPayload, err = authenticatedClient.GenerateDownloadPayload("vmware_tools", "vmtools", "11.1.1", "VMware-Tools-darwin-*.tar.gz", true)
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
	downloadPayload, err = authenticatedClient.GenerateDownloadPayload("vmware_tools", "vmtools", "666", "VMware-Tools-darwin-*.tar.gz", true)
	assert.ErrorIs(t, err, ErrorInvalidVersion)
	assert.Empty(t, downloadPayload, "Expected response to be empty")
}

func TestFetchDownloadLinkNeedEula(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadPayload []DownloadPayload
	downloadPayload, err = authenticatedClient.GenerateDownloadPayload("vmware_tools", "vmtools", "11.1.0", "VMware-Tools-darwin-*.tar.gz", false)
	assert.ErrorIs(t, err, ErrorEulaUnaccepted)
	assert.Empty(t, downloadPayload, "Expected response to be empty")
}

func TestFetchDownloadLinkNotEntitled(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadPayload []DownloadPayload
	downloadPayload, err = authenticatedClient.GenerateDownloadPayload("vmware_nsx_t_data_center", "nsx-t", "3.1.3", "nsx-unified-appliance-secondary-*.qcow2", true)
	assert.ErrorIs(t, err, ErrorNotEntitled)
	assert.Empty(t, downloadPayload, "Expected response to be empty")
}

func TestGenerateDownloadInvalidVersionGlob(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadPayload []DownloadPayload
	downloadPayload, err = authenticatedClient.GenerateDownloadPayload("vmware_tools", "vmtools", "666.*", "VMware-Tools-darwin-*.tar.gz", true)
	assert.ErrorIs(t, err, ErrorNoMatchingVersions)
	assert.Empty(t, downloadPayload, "Expected response to be empty")
}

func TestGenerateDownloadDoubleVersion(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var downloadPayload []DownloadPayload
	downloadPayload, err = authenticatedClient.GenerateDownloadPayload("vmware_tools", "vmtools", "*.*", "VMware-Tools-darwin-*.tar.gz", true)
	assert.ErrorIs(t, err, ErrorMultipleVersionGlob)
	assert.Empty(t, downloadPayload, "Expected response to be empty")
}
