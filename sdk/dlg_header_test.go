package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHeaderSuccess(t *testing.T) {
	var dlgHeader DlgHeader
	dlgHeader, err = basicClient.GetDlgHeader("VMTOOLS1130", "1073")
	assert.Nil(t, err)
	assert.Equal(t, dlgHeader.Product.Productmap, "vmware_tools", "Expected product name vmware_tools")
}

func TestGetHeaderInvalidProductId(t *testing.T) {
	var dlgHeader DlgHeader
	dlgHeader, err = basicClient.GetDlgHeader("VMTOOLS1130", "666666")
	assert.ErrorIs(t, err, ErrorDlgHeader)
	assert.Empty(t, dlgHeader, "Expected response to be empty")
}

func TestGetHeaderInvalidDownloadGroup(t *testing.T) {
	var dlgHeader DlgHeader
	dlgHeader, err = basicClient.GetDlgHeader("VMTOOLS666", "1073")
	assert.ErrorIs(t, err, ErrorDlgHeader)
	assert.Empty(t, dlgHeader, "Expected response to be empty")
}
