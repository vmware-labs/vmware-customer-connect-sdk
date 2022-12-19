package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
)

type DlgDetails struct {
	DownloadDetails     []DownloadDetails   `json:"downloadFiles"`
	EligibilityResponse EligibilityResponse `json:"eligibilityResponse"`
	EulaResponse        EulaResponse        `json:"eulaResponse"`
}

type DownloadDetails struct {
	FileName       string `json:"fileName"`
	Sha1Checksum   string `json:"sha1checksum"`
	Sha256Checksum string `json:"sha256checksum"`
	Md5Checksum    string `json:"md5checksum"`
	Build          string `json:"build"`
	ReleaseDate    string `json:"releaseDate"`
	FileType       string `json:"fileType"`
	Description    string `json:"description"`
	FileSize       string `json:"fileSize"`
	Title          string `json:"title"`
	Version        string `json:"version"`
	Status         string `json:"status"`
	UUID           string `json:"uuid"`
	Header         bool   `json:"header"`
	DisplayOrder   int    `json:"displayOrder"`
	Relink         bool   `json:"relink"`
	Rsync          bool   `json:"rsync"`
}

type EligibilityResponse struct {
	EligibleToDownload bool `json:"eligibleToDownload"`
}
type EulaResponse struct {
	EulaAccepted bool   `json:"eulaAccepted"`
	EulaURL      string `json:"eulaURL"`
}

type FoundDownload struct {
	DownloadDetails    []DownloadDetails
	EulaAccepted       bool
	EligibleToDownload bool
}

const (
	dlgDetailsURLAuthenticated = baseURL + "/channel/api/v1.0/dlg/details"
	dlgDetailsURLPublic        = baseURL + "/channel/public/api/v1.0/dlg/details"
)

var (
	ErrorDlgDetailsInputs      = errors.New("dlgDetails: downloadGroup or productId invalid")
	ErrorNoMatchingFiles       = errors.New("dlgDetails: no files match provided glob")
	ErrorMultipleMatchingFiles = errors.New("dlgDetails: more than 1 file matches glob")
	ErrorEulaUnaccepted        = errors.New("dlgDetails: EULA needs to be accepted for this version")
	ErrorNotEntitled           = errors.New("dlgDetails: user is not entitled to download this file")
)

// curl "https://my.vmware.com/channel/public/api/v1.0/dlg/details?downloadGroup=VMTOOLS1130&productId=1073" |jq
func (c *Client) GetDlgDetails(downloadGroup, productId string) (data DlgDetails, err error) {
	err = c.CheckLoggedIn()
	// Use public URL when user is not logged in
	// This will not return entitlement or EULA sections
	var dlgDetailsURL string
	if err != nil {
		dlgDetailsURL = dlgDetailsURLPublic
	} else {
		dlgDetailsURL = dlgDetailsURLAuthenticated
	}

	search_string := fmt.Sprintf("?downloadGroup=%s&productId=%s", downloadGroup, productId)
	var res *http.Response
	res, err = c.HttpClient.Get(dlgDetailsURL + search_string)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode == 400 {
		err = ErrorDlgDetailsInputs
		return
	} else if res.StatusCode == 401 {
		err = ErrorNotAuthenticated
		return
	}

	err = json.NewDecoder(res.Body).Decode(&data)

	return
}

func (c *Client) FindDlgDetails(downloadGroup, productId, fileName string) (data FoundDownload, err error) {
	if err = c.CheckLoggedIn(); err != nil {
		return
	}

	var dlgDetails DlgDetails
	dlgDetails, err = c.GetDlgDetails(downloadGroup, productId)
	if err != nil {
		return
	}

	data = FoundDownload{
		EulaAccepted:       dlgDetails.EulaResponse.EulaAccepted,
		EligibleToDownload: dlgDetails.EligibilityResponse.EligibleToDownload,
	}

	// Search for file which matches the pattern. If glob is used multiple will return.
	for _, download := range dlgDetails.DownloadDetails {
		filename := download.FileName
		if match, _ := filepath.Match(fileName, filename); match {

			data.DownloadDetails = append(data.DownloadDetails, download)
		}
	}

	if len(data.DownloadDetails) == 0 {
		err = ErrorNoMatchingFiles
	}
	return
}

func (c *Client) GetFileArray(slug, subProduct, version string) (data []string, err error) {
	var downloadGroup, productID string
	downloadGroup, productID, err = c.GetDlgProduct(slug, subProduct, version)
	if err != nil {
		return
	}

	var dlgDetails DlgDetails
	dlgDetails, err = c.GetDlgDetails(downloadGroup, productID)
	if err != nil {
		return
	}

	for _, download := range dlgDetails.DownloadDetails {
		if download.FileName != "" {
			data = append(data, download.FileName)
		}
	}

	return
}

func (c *Client) GetDlgProduct(slug, subProduct, version string) (downloadGroup, productID string, err error) {
	// Find the API version details
	var apiVersion APIVersions
	apiVersion, err = c.FindVersion(slug, subProduct, version)
	if err != nil {
		return
	}

	var subProductDetails DlgList
	subProductDetails, err = c.GetSubProductDetails(slug, subProduct, apiVersion.MajorVersion)
	if err != nil {
		return
	}

	downloadGroup = apiVersion.Code
	productID = subProductDetails.ProductID

	return
}
