package sdk

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/orirawlings/persistent-cookiejar"
	"github.com/stretchr/testify/assert"
)

var authenticatedClient *Client

func ensureLogin(t *testing.T) (err error) {
	if authenticatedClient == nil {
		var jar *cookiejar.Jar
		// Persist cookies on file system to speed up testing
		jar, err = cookiejar.New(&cookiejar.Options{
			Filename:              filepath.Join(homeDir(), ".vmware.cookies"),
			PersistSessionCookies: true,
		})
		if err != nil {return}
		user, pass := mustEnv(t, "VMWCC_USER"), mustEnv(t, "VMWCC_PASS")
		authenticatedClient, err = Login(user, pass, jar)
		if err == nil {
			err = jar.Save()
		}
	}
	return
}

// homeDir returns the OS-specific home path as specified in the environment.
func homeDir() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"))
	}
	return os.Getenv("HOME")
}

func TestSuccessfulLogin(t *testing.T) {
	jar, _ := cookiejar.New(&cookiejar.Options{NoPersist: true})
	user, pass := mustEnv(t, "VMWCC_USER"), mustEnv(t, "VMWCC_PASS")
	_, err = Login(user, pass, jar)
	assert.Nil(t, err)
}

func TestFailedLogin(t *testing.T) {
	jar, _ := cookiejar.New(&cookiejar.Options{NoPersist: true})
	_, err = Login("user", "pass", jar)
	assert.ErrorIs(t, err, ErrorAuthenticationFailure)
}

func TestSuccessfulConnection(t *testing.T) {
	err = CheckConnectivity()
	if err != nil {
		t.Errorf("Expected error not to occur, got %q", err)
	}
}

// func TestInvalidProxy(t *testing.T) {
// 	os.Setenv("HTTPS_PROXY", "http://NOT_A_PROXY")
// 	defer os.Unsetenv("HTTPS_PROXY")
// 	err := CheckConnectivity()
// 	if err == nil {
// 		t.Errorf("Did not generate error when invalid proxy set")
// 	}
// }
