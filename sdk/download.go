package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type DownloadPayload struct {
	Locale        string `json:"locale"`        // en_US
	DownloadGroup string `json:"downloadGroup"` // Versions versionMap[version].Code
	ProductId     string `json:"productId"`     // DlgList ProductID
	Md5checksum   string `json:"md5checksum"`   // dlgDetails  Md5Checksum
	TagId         int    `json:"tagId"`         // dlgHeader Dlg.TagID
	UUId          string `json:"uUId"`          // dlgDetails UUID
	DlgType       string `json:"dlgType"`       // dlgHeader Dlg.Type replace(/&amp;/g, '&')
	ProductFamily string `json:"productFamily"` // dlgHeader Product.Name
	ReleaseDate   string `json:"releaseDate"`   // dlgDetails ReleaseDate
	DlgVersion    string `json:"dlgVersion"`    // dlgDetails Version
	IsBetaFlow    bool   `json:"isBetaFlow"`    // false
}

type AuthorizedDownload struct {
	DownloadURL string `json:"downloadURL"`
	FileName    string `json:"fileName"`
}

const (
	downloadURL = baseURL + "/channel/api/v1.0/dlg/download"
)

var ErrorInvalidDownloadPayload = errors.New("download: invalid download payload")

func (c *Client) GenerateDownloadPayload(slug, subProduct, version, fileName string, acceptEula bool) (data []DownloadPayload, err error) {
	if err = c.CheckLoggedIn(); err != nil {
		return
	}

	if err = c.EnsureProductDetailMap(); err != nil {
		return
	}

	if _, ok := ProductDetailMap[slug]; !ok {
		err = ErrorInvalidSlug
		return
	}

	var downloadGroup, productID string
	downloadGroup, productID, err = c.GetDlgProduct(slug, subProduct, version)
	if err != nil {
		return
	}

	var dlgHeader DlgHeader
	dlgHeader, err = c.GetDlgHeader(downloadGroup, productID)
	if err != nil {
		return
	}

	var downloadDetails FoundDownload
	downloadDetails, err = c.FindDlgDetails(downloadGroup, productID, fileName)
	if err != nil {
		return
	}

	if !downloadDetails.EligibleToDownload {
		err = ErrorNotEntitled
		return
	}

	if !downloadDetails.EulaAccepted {
		if !acceptEula {
			err = ErrorEulaUnaccepted
			return
		} else {
			err = c.AcceptEula(downloadGroup, productID)
			if err != nil {
				return
			}
		}
	}

	for _, downloadFile := range downloadDetails.DownloadDetails {
		downloadPayload := DownloadPayload{
			Locale:        "en_US",
			DownloadGroup: downloadGroup,
			ProductId:     productID,
			Md5checksum:   downloadFile.Md5Checksum,
			TagId:         dlgHeader.Dlg.TagID,
			UUId:          downloadFile.UUID,
			DlgType:       dlgHeader.Dlg.Type,
			ProductFamily: dlgHeader.Product.Name,
			ReleaseDate:   downloadFile.ReleaseDate,
			DlgVersion:    downloadFile.Version,
			IsBetaFlow:    false,
		}

		data = append(data, downloadPayload)

	}

	return
}

func (c *Client) FetchDownloadLink(downloadPayload DownloadPayload) (data AuthorizedDownload, err error) {
	if err = c.CheckLoggedIn(); err != nil {
		return
	}

	postJson, _ := json.Marshal(downloadPayload)
	payload := bytes.NewBuffer(postJson)

	var req *http.Request
	req, err = http.NewRequest("POST", downloadURL, payload)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-XSRF-TOKEN", c.XsrfToken)
	var res *http.Response
	res, err = c.HttpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		err = json.NewDecoder(res.Body).Decode(&data)
	} else if res.StatusCode == 400 {
		err = ErrorInvalidDownloadPayload
	} else {
		err = ErrorNon200Response
	}

	return
}
