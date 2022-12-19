package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDlgList(t *testing.T) {
	var dlgEditions []DlgEditionsLists
	dlgEditions, err = basicClient.GetDlgEditionsList("vmware_vsphere", "7_0")
	assert.Nil(t, err)
	assert.Greater(t, len(dlgEditions), 1, "Expected response to contain more that one item")
}

func TestGetDlgListInvalidSlug(t *testing.T) {
	var dlgEditions []DlgEditionsLists
	dlgEditions, err = basicClient.GetDlgEditionsList("mware_tools", "11_x")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrorInvalidSlug)
	assert.Empty(t, dlgEditions, "Expected response to be empty")
}

func TestGetDlgListInvalidVersion(t *testing.T) {
	var dlgEditions []DlgEditionsLists
	dlgEditions, err = basicClient.GetDlgEditionsList("vmware_tools", "99_x")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrorInvalidVersion)
	assert.Empty(t, dlgEditions, "Expected response to be empty")
}
