// Command gen-product-map is the helper tool to generate lifxlan ProductMap.
//
// To install it:
//
//     go get -u go.yhsif.com/lifxlan/cmd/gen-product-map
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
	"sort"

	"go.yhsif.com/lifxlan"
)

var vendors []struct {
	ID       uint32           `json:"vid"`
	Name     string           `json:"name"`
	Defaults lifxlan.Features `json:"defaults"`

	Products []lifxlan.Product `json:"products"`
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

	err = json.NewDecoder(resp.Body).Decode(&vendors)
	if err != nil {
		errorOut(err)
	}

	fmt.Println(`var ProductMap = map[uint64]Product{`)

	for _, vendor := range vendors {
		for _, product := range vendor.Products {
			product.Features = lifxlan.MergeFeatures(product.Features, vendor.Defaults)
			sort.Sort(sort.Reverse(product.Upgrades))

			fmt.Printf("\tProductMapKey(%v, %v): {\n", vendor.ID, product.ProductID)
			fmt.Printf("\t\tVendorName:  %q,\n", vendor.Name)
			fmt.Printf("\t\tVendorID:    %v,\n", vendor.ID)
			fmt.Printf("\t\tProductName: %q,\n", product.ProductName)
			fmt.Printf("\t\tProductID:   %v,\n", product.ProductID)
			fmt.Printf("\t\tFeatures: Features{\n")
			fmt.Printf("\t\t\tHEV:               OptionalBoolPtr(%v),\n", product.Features.HEV.Get())
			fmt.Printf("\t\t\tColor:             OptionalBoolPtr(%v),\n", product.Features.Color.Get())
			fmt.Printf("\t\t\tChain:             OptionalBoolPtr(%v),\n", product.Features.Chain.Get())
			fmt.Printf("\t\t\tMatrix:            OptionalBoolPtr(%v),\n", product.Features.Matrix.Get())
			fmt.Printf("\t\t\tRelays:            OptionalBoolPtr(%v),\n", product.Features.Relays.Get())
			fmt.Printf("\t\t\tButtons:           OptionalBoolPtr(%v),\n", product.Features.Buttons.Get())
			fmt.Printf("\t\t\tInfrared:          OptionalBoolPtr(%v),\n", product.Features.Infrared.Get())
			fmt.Printf("\t\t\tMultizone:         OptionalBoolPtr(%v),\n", product.Features.Multizone.Get())
			fmt.Printf("\t\t\tExtendedMultizone: OptionalBoolPtr(%v),\n", product.Features.ExtendedMultizone.Get())
			if product.Features.TemperatureRange.Valid() {
				fmt.Printf(
					"\t\t\tTemperatureRange:  TemperatureRange{%v, %v},\n",
					product.Features.TemperatureRange.Min(),
					product.Features.TemperatureRange.Max(),
				)
			}
			fmt.Printf("\t\t},\n")
			if len(product.Upgrades) > 0 {
				fmt.Printf("\t\tUpgrades: Upgrades{\n")
				for _, upgrade := range product.Upgrades {
					fmt.Printf("\t\t\t{\n")
					fmt.Printf("\t\t\t\tMajor: %v,\n", upgrade.Major)
					fmt.Printf("\t\t\t\tMinor: %v,\n", upgrade.Minor)
					fmt.Printf("\t\t\t\tFeatures: Features{\n")
					if upgrade.Features.HEV != nil {
						fmt.Printf("\t\t\t\t\tHEV: OptionalBoolPtr(%v),\n", upgrade.Features.HEV.Get())
					}
					if upgrade.Features.Color != nil {
						fmt.Printf("\t\t\t\t\tColor: OptionalBoolPtr(%v),\n", upgrade.Features.Color.Get())
					}
					if upgrade.Features.Chain != nil {
						fmt.Printf("\t\t\t\t\tChain: OptionalBoolPtr(%v),\n", upgrade.Features.Chain.Get())
					}
					if upgrade.Features.Matrix != nil {
						fmt.Printf("\t\t\t\t\tMatrix: OptionalBoolPtr(%v),\n", upgrade.Features.Matrix.Get())
					}
					if upgrade.Features.Relays != nil {
						fmt.Printf("\t\t\t\t\tRelays: OptionalBoolPtr(%v),\n", upgrade.Features.Relays.Get())
					}
					if upgrade.Features.Buttons != nil {
						fmt.Printf("\t\t\t\t\tButtons: OptionalBoolPtr(%v),\n", upgrade.Features.Buttons.Get())
					}
					if upgrade.Features.Infrared != nil {
						fmt.Printf("\t\t\t\t\tInfrared: OptionalBoolPtr(%v),\n", upgrade.Features.Infrared.Get())
					}
					if upgrade.Features.Multizone != nil {
						fmt.Printf("\t\t\t\t\tMultizone: OptionalBoolPtr(%v),\n", upgrade.Features.Multizone.Get())
					}
					if upgrade.Features.ExtendedMultizone != nil {
						fmt.Printf("\t\t\t\t\tExtendedMultizone: OptionalBoolPtr(%v),\n", upgrade.Features.ExtendedMultizone.Get())
					}
					if upgrade.Features.TemperatureRange.Valid() {
						fmt.Printf(
							"\t\t\t\t\tTemperatureRange: TemperatureRange{%v, %v},\n",
							upgrade.Features.TemperatureRange.Min(),
							upgrade.Features.TemperatureRange.Max(),
						)
					}
					fmt.Printf("\t\t\t\t},\n")
					fmt.Printf("\t\t\t},\n")
				}
				fmt.Printf("\t\t},\n")
			}
			fmt.Printf("\t},\n")
		}
	}

	fmt.Println(`}`)
}
