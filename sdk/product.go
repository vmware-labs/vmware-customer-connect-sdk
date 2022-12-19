package sdk

import (
	"encoding/json"
	"net/http"
	"strings"
)

var ProductDetailMap map[string]ProductDetails

type ProductDetails struct {
	Category           string
	DisplayName        string
	LatestMajorVersion string
}
type ProductResponse struct {
	ProductCategoryList []ProductCategoryList `json:"productCategoryList"`
}
type MajorProductEntities struct {
	Linkname string `json:"linkname"`
	OrderID  int    `json:"orderId"`
	Target   string `json:"target"`
}
type MajorProducts struct {
	MajorProductEntities []MajorProductEntities `json:"actions"`
	Name                 string                 `json:"name"`
}
type ProductCategoryList struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	MajorProducts []MajorProducts `json:"productList"`
}

func (c *Client) GetProductsSlice() (data []MajorProducts, err error) {
	var res *http.Response
	res, err = c.HttpClient.Get(baseURL + "/channel/public/api/v1.0/products/getProductsAtoZ?isPrivate=true")
	if err != nil {
		return
	}
	defer res.Body.Close()

	var decodedProducts ProductResponse
	err = json.NewDecoder(res.Body).Decode(&decodedProducts)

	if err == nil {
		data = decodedProducts.ProductCategoryList[0].MajorProducts
	}

	return
}

// returned map is used to look up products by their slig
func (c *Client) GetProductsMap() (productMap map[string]ProductDetails, err error) {
	productMap = make(map[string]ProductDetails)

	var products []MajorProducts
	products, err = c.GetProductsSlice()
	if err != nil {
		return
	}

	for _, product := range products {
		for _, subProduct := range product.MajorProductEntities {
			if !strings.Contains(subProduct.Target, "http") {
				splitTarget := strings.Split(subProduct.Target, "/")
				productDetails := ProductDetails{
					Category:           splitTarget[3],
					DisplayName:        product.Name,
					LatestMajorVersion: strings.Split(splitTarget[5], "#")[0],
				}
				productMap[splitTarget[4]] = productDetails
			}
		}
	}
	return
}

func (c *Client) EnsureProductDetailMap() (err error) {
	if len(ProductDetailMap) < 1 {
		ProductDetailMap, err = c.GetProductsMap()
	}
	return
}

func (c *Client) GetCategory(slug string) (data string, err error) {
	c.EnsureProductDetailMap()
	if err != nil {
		return
	}

	if category, ok := ProductDetailMap[slug]; ok {
		data = category.Category
	} else {
		err = ErrorInvalidSlug
	}

	return
}
