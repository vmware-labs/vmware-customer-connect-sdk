// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache 2.0

package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDlgListDownloads(t *testing.T) {
	var dlgEditions []DlgEditionsLists
	dlgEditions, err = basicClient.GetDlgEditionsList("vmware_vsphere", "7_0", "PRODUCT_BINARY")
	assert.Nil(t, err)
	assert.Greater(t, len(dlgEditions), 1, "Expected response to contain more that one item")
}

func TestGetDlgListDrivers(t *testing.T) {
	var dlgEditions []DlgEditionsLists
	dlgEditions, err = basicClient.GetDlgEditionsList("vmware_vsphere", "7_0", "DRIVERS_TOOLS")
	assert.Nil(t, err)
	assert.Greater(t, len(dlgEditions), 1, "Expected response to contain more that one item")
}

func TestGetDlgListCustomISO(t *testing.T) {
	var dlgEditions []DlgEditionsLists
	dlgEditions, err = basicClient.GetDlgEditionsList("vmware_vsphere", "7_0", "CUSTOM_ISO")
	assert.Nil(t, err)
	assert.Greater(t, len(dlgEditions), 0, "Expected response to contain more that one item")
}

func TestGetDlgListAddons(t *testing.T) {
	var dlgEditions []DlgEditionsLists
	dlgEditions, err = basicClient.GetDlgEditionsList("vmware_vsphere", "7_0", "ADDONS")
	assert.Nil(t, err)
	assert.Greater(t, len(dlgEditions), 0, "Expected response to contain more that one item")
}

func TestGetDlgListInvalidSlug(t *testing.T) {
	var dlgEditions []DlgEditionsLists
	dlgEditions, err = basicClient.GetDlgEditionsList("mware_tools", "11_x", "PRODUCT_BINARY")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrorInvalidSlug)
	assert.Empty(t, dlgEditions, "Expected response to be empty")
}

func TestGetDlgListInvalidVersion(t *testing.T) {
	var dlgEditions []DlgEditionsLists
	dlgEditions, err = basicClient.GetDlgEditionsList("vmware_tools", "99_x", "PRODUCT_BINARY")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrorInvalidVersion)
	assert.Empty(t, dlgEditions, "Expected response to be empty")
}
