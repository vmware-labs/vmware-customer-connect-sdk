// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache 2.0

package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchEulaLink(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var eulaUrl string
	eulaUrl, err = authenticatedClient.FetchEulaUrl("VMTOOLS1235", "1259")
	assert.Nil(t, err)
	assert.NotEmpty(t, eulaUrl, "Expected eulaUrl not be empty")
}

func TestFetchEulaLinkInvalidCode(t *testing.T) {
	err = ensureLogin(t)
	assert.Nil(t, err)

	var eulaUrl string
	eulaUrl, err = authenticatedClient.FetchEulaUrl("VMTOOLS1235", "9999")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrorDlgDetailsInputs)
	assert.Empty(t, eulaUrl, "Expected eulaUrl be empty")
}

func TestAcceptEula(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	err = authenticatedClient.AcceptEula("VMTOOLS1235", "1259")
	assert.Nil(t, err)
}
