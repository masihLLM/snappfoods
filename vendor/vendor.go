package vendor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Vendor struct {
	Code string `json:"code"`
}

type Data struct {
	Data Vendor `json:"data"`
}

type Response struct {
	Status bool `json:"status"`
	Data   struct {
		Vendor interface{} `json:"vendor"`
	} `json:"data"`
}

func ReadVendorCodes(root string) {
	vendorCodes := []string{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".json" {
			file, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			var data Data
			if err := json.Unmarshal(file, &data); err != nil {
				return err
			}

			vendorCodes = append(vendorCodes, data.Data.Code)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to read files: %v", err)
	}

	for _, code := range vendorCodes {
		url := fmt.Sprintf("https://snappfood.ir/mobile/v2/restaurant/details/dynamic?lat=-1&long=-1&optionalClient=WEBSITE&client=WEBSITE&deviceType=WEBSITE&appVersion=8.1.1&UDID=1351f4cb-a3c7-4033-995e-31776b068f93&vendorCode=%s&locationCacheKey=lat%%3D-1%%26long%%3D-1&show_party=1&fetch-static-data=1&locale=fa", code)

		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response body: %v", err)
		}

		var response Response
		if err := json.Unmarshal(body, &response); err != nil {
			log.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Status {
			filePath := filepath.Join("vendor", fmt.Sprintf("%s.json", code))
			if err := ioutil.WriteFile(filePath, body, 0644); err != nil {
				log.Fatalf("Failed to write file: %v", err)
			}
		}
	}
}
