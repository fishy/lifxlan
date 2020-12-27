package light_test

import (
	"fmt"
	"testing"

	"go.yhsif.com/lifxlan/light"
)

func TestBool2Uint8(t *testing.T) {
	cases := map[bool]light.BoolUint8{
		false: 0,
		true:  1,
	}

	for v, expected := range cases {
		t.Run(
			fmt.Sprintf("%v", v),
			func(t *testing.T) {
				got := light.Bool2Uint8(v)
				if got != expected {
					t.Errorf("Bool2Uint8(%v) expected %d, got %d", v, expected, got)
				}
			},
		)
	}
}

func TestConvertSkewRatio(t *testing.T) {
	cases := map[float64]int16{
		0:    -32768,
		0.25: -16384,
		0.5:  0,
		0.75: 16383,
		1:    32767,
	}

	for v, expected := range cases {
		t.Run(
			fmt.Sprintf("%v", v),
			func(t *testing.T) {
				got := light.ConvertSkewRatio(v)
				if got != expected {
					t.Errorf("ConvertSkewRatio(%v) expected %d, got %d", v, expected, got)
				}
			},
		)
	}
}
