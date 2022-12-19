package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMajorVersionsSuccess(t *testing.T) {
	var majorVersions []string
	majorVersions, err = basicClient.GetMajorVersionsSlice("vmware_tools")
	assert.Nil(t, err)
	assert.Greater(t, len(majorVersions), 1, "Expected response to contain at least 1 item")
	assert.Contains(t, majorVersions, "11_x")
}

func TestGetMajorVersionsInvalidSlug(t *testing.T) {
	var majorVersions []string
	majorVersions, err = basicClient.GetMajorVersionsSlice("mware_tools")
	assert.ErrorIs(t, err, ErrorInvalidSlug)
	assert.Empty(t, majorVersions, "Expected response to be empty")
}
