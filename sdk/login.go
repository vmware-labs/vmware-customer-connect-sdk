package sdk

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/orirawlings/persistent-cookiejar"
	"golang.org/x/net/html"
)

type Client struct {
	HttpClient *http.Client
	XsrfToken  string
}

type TokenValidation struct {
	IsAccessTokenValid    string `json:"isAccessTokenValid"`
	IsRefreshTokenPresent string `json:"isRefreshTokenPresent"`
	Cn                    string `json:"cn"`
	Email                 string `json:"email"`
	Firstname             string `json:"firstname"`
	Lastname              string `json:"lastname"`
	UserType              string `json:"userType"`
}

const (
	baseURL = "https://customerconnect.vmware.com"
	initURL = baseURL + "/web/vmware/login"
	ssoURL  = baseURL + "/vmwauth/saml/SSO"
	authURL = "https://auth.vmware.com/oam/server/auth_cred_submit?Auth-AppID=WMVMWR"
)

const samlInputQuery = `input[name="SAMLResponse"]`

var samlInputSelector = cascadia.MustCompile(samlInputQuery)

var ErrorNotAuthenticated = errors.New("generic: returned http 401 not authenticated")
var ErrorAuthenticationFailure = errors.New("login: authentication failure")
var ErrorXsrfFailure = errors.New("login: server did not return XSRF token")
var ErrorConnectionFailure = errors.New("login: server did not return 200 ok")

func Login(username, password string, jar *cookiejar.Jar) (client *Client, err error) {
	err = CheckConnectivity()
	if err != nil {
		return
	}

	httpClient := &http.Client{Jar: jar}

	// When cookies are passed in and check to see can make calls
	// Otherwise perform a login
	_, errXsrf := setXsrfToken(httpClient)
	loginNeeded := false
	if len(jar.AllCookies()) > 0 && errXsrf == nil {

		payload := `{"rowLimit": 10}`
		var res *http.Response
		res, err = httpClient.Post(accountInfoURL, "application/json", strings.NewReader(payload))
		if err != nil {
			return
		}
		defer res.Body.Close()

		if res.StatusCode == 401 || res.StatusCode == 500 {
			loginNeeded = true
		}
	} else {
		loginNeeded = true
	}

	if loginNeeded {
		jar.RemoveAll()
		err = performLogin(httpClient, username, password)
		if err != nil {
			return
		}
	}

	var xsrfToken string
	if xsrfToken, err = setXsrfToken(httpClient); err != nil {
		return
	}

	client = &Client{
		HttpClient: httpClient,
		XsrfToken:  xsrfToken,
	}

	return
}

// Extract xsrf token value to be used when getting download link
func setXsrfToken(client *http.Client) (xsrfToken string, err error) {
	u, _ := url.Parse("https://customerconnect.vmware.com")
	cookies := client.Jar.Cookies(u)
	for _, cookie := range cookies {
		if cookie.Name == "XSRF-TOKEN" {
			xsrfToken = cookie.Value
		}
	}
	if xsrfToken == "" {
		err = ErrorXsrfFailure
		return
	}
	return
}

func performLogin(httpClient *http.Client, username, password string) (err error) {

	// Initialize cookies
	initRes, err := httpClient.Get(initURL)
	if err != nil {
		return
	}
	initRes.Body.Close()
	if initRes.StatusCode != 200 {
		err = ErrorConnectionFailure
		return
	}

	var authResp *http.Response
	var buf bytes.Buffer

	// Attempting login 5 times, to allow for empty response issue.
	for i := 1; i < 5; i++ {
		// Post credentials to get SAML token back
		authResp, err = httpClient.PostForm(authURL, url.Values{
			"username": {username},
			"password": {password},
		})
		if err != nil {
			break
		}
		defer authResp.Body.Close()

		tee := io.TeeReader(authResp.Body, &buf)

		var authBodyBytes []byte
		authBodyBytes, err = io.ReadAll(tee)
		str_body := string(authBodyBytes)

		if str_body == "" {
			continue
		}

		if authResp.Request.URL.Path == "/login" {
			// Return auth failure if reposonse is redirect to the login page
			err = ErrorAuthenticationFailure
		}
		break
	}

	if err != nil {
		return
	}

	samlToken, err := getSAMLToken(&buf)
	if err != nil {
		return
	}

	// Post SAML token to generate final session cookies
	ssoRes, err := httpClient.PostForm(ssoURL, url.Values{
		"SAMLResponse": {samlToken},
	})
	if err != nil {
		return
	}
	defer ssoRes.Body.Close()

	return
}

// Extract SAML token from HTML body
func getSAMLToken(body io.Reader) (string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return "", err
	}

	node := samlInputSelector.MatchFirst(doc)
	if node == nil {
		return "", fmt.Errorf("could not find node that matches %#v", samlInputQuery)
	}

	for _, attr := range node.Attr {
		if attr.Key == "value" {
			return attr.Val, nil
		}
	}

	return "", fmt.Errorf("could not find the node's value attribute")
}

func CheckConnectivity() (err error) {
	httpClient := &http.Client{}
	var res *http.Response
	res, err = httpClient.Get(ssoURL)

	if res.StatusCode != 200 {
		err = ErrorConnectionFailure
	}

	return
}
