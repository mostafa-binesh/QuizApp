package services

import (
	U "docker/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// import (
// 	U "docker/utils"
// 	"github.com/chenyangguang/woocommerce"
// )

type WCProduct struct {
	ID         int                   `json:"id"`
	Name       string                `json:"name"`
	Attributes []WCProductAttributes `json:"attributes"`
}
type WCProductAttributes struct {
	Name    string   `json:"name"`
	Options []string `json:"options"`
}

// in wc. products model, we have serveral attributes:
// 1. parent courses (products) have only duration attr.
// 2. non-parent courses have ValidityPeriod and ParentCourse (parentCourseID)

func GetAllProducts() ([]WCProduct, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/products?per_page=100", U.Env("WC_BASE_URL"))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(U.Env("WC_CONSUMER_KEY"), U.Env("WC_CONSUMER_SECRET"))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve products. status code: %d", resp.StatusCode)
	}
	var products []WCProduct
	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func GetProductByID(productID int) (WCProduct, error) {
	client := &http.Client{}

	url := fmt.Sprintf("%s/products/%d", U.Env("WC_BASE_URL"), productID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return WCProduct{}, err
	}

	req.SetBasicAuth(U.Env("WC_CONSUMER_KEY"), U.Env("WC_CONSUMER_SECRET"))

	resp, err := client.Do(req)
	if err != nil {
		return WCProduct{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WCProduct{}, fmt.Errorf("failed to retrieve product. status code: %d", resp.StatusCode)
	}

	var product WCProduct
	err = json.NewDecoder(resp.Body).Decode(&product)
	if err != nil {
		return WCProduct{}, err
	}

	return product, nil
}
