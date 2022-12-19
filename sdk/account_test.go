package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountInfo(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)

	var accountInfo AccountInfo
	accountInfo, err = authenticatedClient.AccountInfo()
	t.Log(accountInfo)
	assert.Nil(t, err)
}

func TestAccountInfoNotLoggedIn(t *testing.T) {
	var accountInfo AccountInfo
	accountInfo, err = basicClient.AccountInfo()
	assert.ErrorIs(t, err, ErrorNotAuthorized)
	assert.Empty(t, accountInfo, "Expected response to be empty")
}

func TestCurrentUser(t *testing.T) {
	err = ensureLogin(t)
	require.Nil(t, err)
	var currentUser CurrentUser
	currentUser, err = authenticatedClient.CurrentUser()
	t.Log(currentUser)
	assert.Nil(t, err)
}

func TestCurrentUserNotLoggedIn(t *testing.T) {
	_, err = basicClient.CurrentUser()
	assert.ErrorIs(t, err, ErrorNotAuthorized)
}

// func TestInvalidProxy(t *testing.T) {
// 	os.Setenv("HTTPS_PROXY", "http://NOT_A_PROXY")
// 	defer os.Unsetenv("HTTPS_PROXY")
// 	err := sdk.CheckConnectivity()
// 	if err == nil {
// 		t.Errorf("Did not generate error when invalid proxy set")
// 	}
// }
