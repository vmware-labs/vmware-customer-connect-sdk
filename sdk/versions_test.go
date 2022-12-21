// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache 2.0

package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersionSuccess(t *testing.T) {
	var versions map[string]APIVersions
	versions, err = basicClient.GetVersionMap("vmware_horizon_clients", "cart+win")
	assert.Nil(t, err)
	assert.Greater(t, len(versions), 1, "Expected response to contain at least 1 item")
	assert.Contains(t, versions, "2106")
	assert.Contains(t, versions, "2006")
}

func TestGetVersionSuccessHorizon(t *testing.T) {
	var versions map[string]APIVersions
	versions, err = basicClient.GetVersionMap("vmware_horizon", "dem+standard")
	assert.Nil(t, err)
	assert.Greater(t, len(versions), 1, "Expected response to contain at least 1 item")
	assert.Contains(t, versions, "2106")
	assert.Contains(t, versions, "2006")
}

func TestGetVersionSuccessNsxLe(t *testing.T) {
	var versions map[string]APIVersions
	versions, err = basicClient.GetVersionMap("vmware_nsx", "nsx_le")
	assert.Nil(t, err)
	assert.Greater(t, len(versions), 1, "Expected response to contain at least 1 item")
	assert.Contains(t, versions, "4.0.1.1 LE")
	assert.NotContains(t, versions, "4.0.1.1")
}

func TestGetVersionSuccessNsx(t *testing.T) {
	var versions map[string]APIVersions
	versions, err = basicClient.GetVersionMap("vmware_nsx", "nsx")
	assert.Nil(t, err)
	assert.Greater(t, len(versions), 1, "Expected response to contain at least 1 item")
	assert.Contains(t, versions, "4.0.1.1")
	assert.NotContains(t, versions, "4.0.1.1 LE")
}

func TestGetVersionMapInvalidSubProduct(t *testing.T) {
	var versions map[string]APIVersions
	versions, err = basicClient.GetVersionMap("vmware_tools", "dummy")
	assert.ErrorIs(t, err, ErrorInvalidSubProduct)
	assert.Empty(t, versions, "Expected response to be empty")
}

func TestGetVersionInvalidSlug(t *testing.T) {
	var versions map[string]APIVersions
	versions, err = basicClient.GetVersionMap("mware_tools", "vmtools")
	assert.ErrorIs(t, err, ErrorInvalidSlug)
	assert.Empty(t, versions, "Expected response to be empty")
}

func TestFindVersion(t *testing.T) {
	var foundVersion APIVersions
	foundVersion, err = basicClient.FindVersion("vmware_tools", "vmtools", "11.1.1")
	assert.Nil(t, err)
	assert.NotEmpty(t, foundVersion.Code, "Expected response not to be empty")
}

func TestFindVersionInvalidSlug(t *testing.T) {
	var foundVersion APIVersions
	foundVersion, err = basicClient.FindVersion("mware_tools", "vmtools", "11.1.1")
	assert.ErrorIs(t, err, ErrorInvalidSlug)
	assert.Empty(t, foundVersion.Code, "Expected response to be empty")
}

func TestFindVersionInvalidVersion(t *testing.T) {
	var foundVersion APIVersions
	foundVersion, err = basicClient.FindVersion("vmware_tools", "vmtools", "666")
	assert.ErrorIs(t, err, ErrorInvalidVersion)
	assert.Empty(t, foundVersion.Code, "Expected response to be empty")
}

func TestFindVersionInvalidSubProduct(t *testing.T) {
	var foundVersion APIVersions
	foundVersion, err = basicClient.FindVersion("vmware_tools", "tools", "11.1.1")
	assert.ErrorIs(t, err, ErrorInvalidSubProduct)
	assert.Empty(t, foundVersion.Code, "Expected response to be empty")
}

func TestFindVersionMinorGlob(t *testing.T) {
	var foundVersion APIVersions
	foundVersion, err = basicClient.FindVersion("vmware_tools", "vmtools", "10.2.*")
	assert.Nil(t, err)
	assert.Equal(t, foundVersion.Code, "VMTOOLS1021")
}

func TestFindVersionOnlyGlob(t *testing.T) {
	var foundVersion APIVersions
	foundVersion, err = basicClient.FindVersion("vmware_tools", "vmtools", "*")
	assert.Nil(t, err)
	assert.NotEmpty(t, foundVersion.Code)
}

func TestGetVersionArraySuccess(t *testing.T) {
	var versions []string
	versions, err = basicClient.GetVersionSlice("vmware_tools", "vmtools")
	assert.Nil(t, err)
	assert.Greater(t, len(versions), 10, "Expected response to contain at least 10 items")
}
