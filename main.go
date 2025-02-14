package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http" 
    "time"
)

type ApiResponse struct {
    Status bool `json:"status"`
    Data   struct {
        FinalResult []interface{} `json:"finalResult"`
    } `json:"data"`
}

func main() {
    counter := 0

    for page := 0; page <= 1000; page++ {
        url := fmt.Sprintf("https://snappfood.ir/search/api/v1/desktop/vendors-list?page=%d&page_size=20&city_name=tehran&locale=fa", page)
        resp, err := http.Get(url)
        if err != nil {
            fmt.Printf("Failed to fetch page %d: %v\n", page, err)
            continue
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            fmt.Printf("Failed to read response body for page %d: %v\n", page, err)
            continue
        }

        var apiResponse ApiResponse
        err = json.Unmarshal(body, &apiResponse)
        if err != nil {
            fmt.Printf("Failed to parse JSON for page %d: %v\n", page, err)
            continue
        }

        if apiResponse.Status {
            for _, result := range apiResponse.Data.FinalResult {
                resultBytes, err := json.Marshal(result)
                if err != nil {
                    fmt.Printf("Failed to marshal result for page %d: %v\n", page, err)
                    continue
                }

                fileName := fmt.Sprintf("%d.json", counter)
                err = ioutil.WriteFile(fileName, resultBytes, 0644)
                if err != nil {
                    fmt.Printf("Failed to write file %s: %v\n", fileName, err)
                    continue
                }

                fmt.Printf("Saved result %d to %s\n", counter, fileName)
                counter++
            }
        } else {
            fmt.Printf("No data found for page %d\n", page)
        }

        time.Sleep(5 * time.Second)
    }
}