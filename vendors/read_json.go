package main

import "x/vendor"

type Data struct {
	Code string `json:"code"`
}

func main() {
	vendor.ReadVendorCodes("../")
}
