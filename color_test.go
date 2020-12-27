package lifxlan_test

import (
	"image/color"
	"reflect"
	"testing"

	"go.yhsif.com/lifxlan"
)

const kelvin = lifxlan.KelvinCool

func TestFromColor(t *testing.T) {
	type testCase struct {
		Label    string
		Color    color.Color
		Expected lifxlan.Color
	}

	cases := []testCase{
		{
			Label: "Black",
			Color: color.Black,
			Expected: lifxlan.Color{
				Hue:        0,
				Saturation: 0,
				Brightness: 0,
			},
		},
		{
			Label: "White",
			Color: color.White,
			Expected: lifxlan.Color{
				Hue:        0,
				Saturation: 0,
				Brightness: 65535,
			},
		},

		{
			Label: "Red",
			Color: &color.RGBA{
				R: 0xff,
			},
			Expected: lifxlan.Color{
				Hue:        0,
				Saturation: 65535,
				Brightness: 65535,
			},
		},
		{
			Label: "Green",
			Color: &color.RGBA{
				G: 0xff,
			},
			Expected: lifxlan.Color{
				Hue:        21845,
				Saturation: 65535,
				Brightness: 65535,
			},
		},
		{
			Label: "Blue",
			Color: &color.RGBA{
				B: 0xff,
			},
			Expected: lifxlan.Color{
				Hue:        43691,
				Saturation: 65535,
				Brightness: 65535,
			},
		},
	}

	for _, test := range cases {
		t.Run(
			test.Label,
			func(t *testing.T) {
				expected := test.Expected
				expected.Kelvin = kelvin
				got := *lifxlan.FromColor(test.Color, kelvin)
				if !reflect.DeepEqual(got, expected) {
					r, g, b, a := test.Color.RGBA()
					t.Errorf(
						"{r:%04x, g:%04x, b:%04x, a:%04x} expected %+v, got %+v",
						r, g, b, a, expected, got,
					)
				}
			},
		)
	}
}
