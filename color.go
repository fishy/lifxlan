package lifxlan

import (
	"image/color"
	"math"
)

// Color is the HSBK color type used in lifx lan API.
//
// https://lan.developer.lifx.com/v2.0/docs/light-messages#section-hsbk
type Color struct {
	Hue        uint16
	Saturation uint16
	Brightness uint16
	Kelvin     uint16
}

// ColorBlack is the black color.
var ColorBlack = *FromColor(color.Black, 0)

// Color value boundaries and constants.
const (
	KelvinWarm uint16 = 2500
	KelvinCool        = 9000

	KelvinMin = KelvinWarm
	KelvinMax = KelvinCool
)

// Sanitize tries to sanitize the color values to keep them within appropriate
// boundaries.
func (c *Color) Sanitize() {
	if c.Kelvin < KelvinMin {
		c.Kelvin = KelvinMin
	}
	if c.Kelvin > KelvinMax {
		c.Kelvin = KelvinMax
	}
}

// FromColor converts a standard library color into HSBK color.
//
// Alpha channel will be ignored and kelvin value will be added.
func FromColor(c color.Color, kelvin uint16) *Color {
	// helper stuff
	const rgbBase = 0xffff
	const hueRate = float64(1<<16) / 360
	const sbRate = float64(math.MaxUint16)
	intMax := func(args ...int) int {
		var max int
		for i, v := range args {
			if i == 0 {
				max = v
				continue
			}
			if v > max {
				max = v
			}
		}
		return max
	}
	intMin := func(args ...int) int {
		var min int
		for i, v := range args {
			if i == 0 {
				min = v
				continue
			}
			if v < min {
				min = v
			}
		}
		return min
	}

	var h, s int

	rr, gg, bb, _ := c.RGBA()
	r := int(rr)
	g := int(gg)
	b := int(bb)

	cmax := intMax(r, g, b)
	delta := cmax - intMin(r, g, b)

	// hue
	switch {
	case delta == 0:
		h = 0
	case cmax == r:
		h = int(math.Round((float64(g-b) / float64(delta)) * 60))
	case cmax == g:
		h = int(math.Round((float64(b-r)/float64(delta) + 2) * 60))
	case cmax == b:
		h = int(math.Round((float64(r-g)/float64(delta) + 4) * 60))
	}
	for h < 0 {
		h += 360
	}
	for h >= 360 {
		h -= 360
	}
	h = int(math.Round(float64(h) * hueRate))

	// saturation
	if cmax == 0 {
		s = 0
	} else {
		s = int(math.Round(float64(delta) / float64(cmax) * sbRate))
	}

	ret := Color{
		Hue:        uint16(h),
		Saturation: uint16(s),
		Brightness: uint16(math.Round(float64(cmax) / rgbBase * sbRate)),
		Kelvin:     kelvin,
	}
	ret.Sanitize()
	return &ret
}
