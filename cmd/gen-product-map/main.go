// Command gen-product-map is the helper tool to generate lifxlan ProductMap.
//
// To install it:
//
//     go get -u github.com/fishy/lifxlan/cmd/gen-product-map
//
// To run it:
//
//     gen-product-map >> product_map.go
//
// Then manally update the file to remove previous value.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

var data []struct {
	ID       uint32 `json:"vid"`
	Name     string `json:"name"`
	Products []struct {
		ID       uint32 `json:"pid"`
		Name     string `json:"name"`
		Features struct {
			Color            bool     `json:"color"`
			Infrared         bool     `json:"infrared"`
			MultiZone        bool     `json:"multizone"`
			Chain            bool     `json:"chain"`
			TemperatureRange []uint16 `json:"temperature_range"`
		} `json:"features"`
	} `json:"products"`
}

var url = flag.String(
	"url",
	"https://raw.githubusercontent.com/LIFX/products/master/products.json",
	"The URL to fetch json data.",
)

func errorOut(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(-1)
}

func main() {
	flag.Parse()
	resp, err := http.Get(*url)
	if err != nil {
		errorOut(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		errorOut(err)
	}

	fmt.Println(`var ProductMap = map[uint64]ParsedHardwareVersion{`)

	for _, vendor := range data {
		for _, product := range vendor.Products {
			if len(product.Features.TemperatureRange) != 2 {
				fmt.Fprintf(os.Stderr, "Skipping invalid product info: %+v\n", product)
				continue
			}
			fmt.Printf("\tProductMapKey(%v, %v): {\n", vendor.ID, product.ID)
			fmt.Printf("\t\tVendorName:  %q,\n", vendor.Name)
			fmt.Printf("\t\tProductName: %q,\n", product.Name)
			fmt.Printf("\t\tColor:       %v,\n", product.Features.Color)
			fmt.Printf("\t\tInfrared:    %v,\n", product.Features.Infrared)
			fmt.Printf("\t\tMultiZone:   %v,\n", product.Features.MultiZone)
			fmt.Printf("\t\tChain:       %v,\n", product.Features.Chain)
			fmt.Printf("\t\tMinKelvin:   %v,\n", product.Features.TemperatureRange[0])
			fmt.Printf("\t\tMaxKelvin:   %v,\n", product.Features.TemperatureRange[1])
			fmt.Printf("\t},\n")
		}
	}

	fmt.Println(`}`)
}
