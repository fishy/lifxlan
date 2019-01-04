package lifxlan_test

import (
	"math/rand"
	"testing"
	"testing/quick"
	"time"

	"github.com/fishy/lifxlan"
)

func TestTime(t *testing.T) {
	// Seed rander
	now := time.Now()
	rander := rand.New(rand.NewSource(now.Unix() + int64(now.Nanosecond())))

	var n int
	var innerT *testing.T

	cases := map[string]func() bool{
		"Timestamp": func() bool {
			n++
			sec := int64(rander.Int31())
			nano := int64(rander.Intn(int(time.Second)))
			t := time.Unix(sec, nano).Round(time.Millisecond)
			ts := lifxlan.ConvertTime(t)
			got := ts.Time()
			if !got.Equal(t) {
				innerT.Logf("Expected %v, got %v", t, got)
				return false
			}
			return true
		},

		"TransitionTime": func() bool {
			n++
			d := time.Duration(rander.Int31()).Round(time.Millisecond)
			ti := lifxlan.ConvertDuration(d)
			got := ti.Duration()
			if got != d {
				innerT.Logf("Expected %v, got %v", d, got)
				return false
			}
			return true
		},
	}

	for label, f := range cases {
		t.Run(
			label,
			func(t *testing.T) {
				n = 0
				innerT = t
				if err := quick.Check(f, nil); err != nil {
					t.Error(err)
				}
				t.Logf("quick did %d checks", n)
			},
		)
	}
}
