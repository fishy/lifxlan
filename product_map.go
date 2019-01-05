package lifxlan

// ProductMap is the map of all known hardwares.
//
// If a new product is added and this file is not updated yet,
// you can add it to the map by yourself, for example:
//
//     func init() {
//         key := lifxlan.ProductMapKey(newVID, newPID)
//         lifxlan.ProductMap[key] = ParsedHardwareVersion{
//				     // Fill in values
//         }
//     }
//
// The content of this map was fetched from
// https://github.com/LIFX/products/blob/master/products.json
// and generated by
// https://github.com/fishy/lifxlan/cmd/gen-product-map/
var ProductMap = map[uint64]ParsedHardwareVersion{
	ProductMapKey(1, 1): {
		ProductName: "Original 1000",
		Color:       true,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 3): {
		ProductName: "Color 650",
		Color:       true,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 10): {
		ProductName: "White 800 (Low Voltage)",
		Color:       false,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2700,
		MaxKelvin:   6500,
	},
	ProductMapKey(1, 11): {
		ProductName: "White 800 (High Voltage)",
		Color:       false,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2700,
		MaxKelvin:   6500,
	},
	ProductMapKey(1, 18): {
		ProductName: "White 900 BR30 (Low Voltage)",
		Color:       false,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2700,
		MaxKelvin:   6500,
	},
	ProductMapKey(1, 20): {
		ProductName: "Color 1000 BR30",
		Color:       true,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 22): {
		ProductName: "Color 1000",
		Color:       true,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 27): {
		ProductName: "LIFX A19",
		Color:       true,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 28): {
		ProductName: "LIFX BR30",
		Color:       true,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 29): {
		ProductName: "LIFX+ A19",
		Color:       true,
		Infrared:    true,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 30): {
		ProductName: "LIFX+ BR30",
		Color:       true,
		Infrared:    true,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 31): {
		ProductName: "LIFX Z",
		Color:       true,
		Infrared:    false,
		MultiZone:   true,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 32): {
		ProductName: "LIFX Z 2",
		Color:       true,
		Infrared:    false,
		MultiZone:   true,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 36): {
		ProductName: "LIFX Downlight",
		Color:       true,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 37): {
		ProductName: "LIFX Downlight",
		Color:       true,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 38): {
		ProductName: "LIFX Beam",
		Color:       true,
		Infrared:    false,
		MultiZone:   true,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 43): {
		ProductName: "LIFX A19",
		Color:       true,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 44): {
		ProductName: "LIFX BR30",
		Color:       true,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 45): {
		ProductName: "LIFX+ A19",
		Color:       true,
		Infrared:    true,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 46): {
		ProductName: "LIFX+ BR30",
		Color:       true,
		Infrared:    true,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 49): {
		ProductName: "LIFX Mini",
		Color:       true,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 50): {
		ProductName: "LIFX Mini Day and Dusk",
		Color:       false,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   1500,
		MaxKelvin:   4000,
	},
	ProductMapKey(1, 51): {
		ProductName: "LIFX Mini White",
		Color:       false,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2700,
		MaxKelvin:   2700,
	},
	ProductMapKey(1, 52): {
		ProductName: "LIFX GU10",
		Color:       true,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 55): {
		ProductName: "LIFX Tile",
		Color:       true,
		Infrared:    false,
		MultiZone:   false,
		Chain:       true,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 56): {
		ProductName: "LIFX Beam",
		Color:       true,
		Infrared:    false,
		MultiZone:   true,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 59): {
		ProductName: "LIFX Mini Color",
		Color:       true,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2500,
		MaxKelvin:   9000,
	},
	ProductMapKey(1, 60): {
		ProductName: "LIFX Mini Day and Dusk",
		Color:       false,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   1500,
		MaxKelvin:   4000,
	},
	ProductMapKey(1, 61): {
		ProductName: "LIFX Mini White",
		Color:       false,
		Infrared:    false,
		MultiZone:   false,
		Chain:       false,
		MinKelvin:   2700,
		MaxKelvin:   2700,
	},
}
