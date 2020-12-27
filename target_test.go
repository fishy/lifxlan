package lifxlan_test

import (
	"flag"
	"fmt"
	"testing"

	"go.yhsif.com/lifxlan"
)

func TestTargetString(t *testing.T) {
	t.Run(
		"AllDevices",
		func(t *testing.T) {
			target := lifxlan.AllDevices
			str := fmt.Sprintf("%v", target)
			expected := "00:00:00:00:00:00"
			if str != expected {
				t.Errorf("Expected %s, got %s", expected, str)
			}
		},
	)

	t.Run(
		"NormalDevice",
		func(t *testing.T) {
			target := lifxlan.Target(0xffffffd573d0)
			str := fmt.Sprintf("%v", target)
			expected := "d0:73:d5:ff:ff:ff"
			if str != expected {
				t.Errorf("Expected %s, got %s", expected, str)
			}
		},
	)
}

func TestTargetMatches(t *testing.T) {
	t.Run(
		"Mismatch",
		func(t *testing.T) {
			t1 := lifxlan.Target(1)
			t2 := lifxlan.Target(2)
			if t1.Matches(t2) {
				t.Errorf("%v should not match %v", t1, t2)
			}
			if t2.Matches(t1) {
				t.Errorf("%v should not match %v", t2, t1)
			}
		},
	)

	t.Run(
		"Same",
		func(t *testing.T) {
			t1 := lifxlan.Target(5)
			t2 := lifxlan.Target(5)
			if !t1.Matches(t2) {
				t.Errorf("%v should match %v", t1, t2)
			}
			if !t2.Matches(t1) {
				t.Errorf("%v should match %v", t2, t1)
			}
		},
	)

	t.Run(
		"AllDevices",
		func(t *testing.T) {
			t1 := lifxlan.Target(1)
			t2 := lifxlan.AllDevices
			if !t1.Matches(t2) {
				t.Errorf("%v should match %v", t1, t2)
			}
			if !t2.Matches(t1) {
				t.Errorf("%v should match %v", t2, t1)
			}
		},
	)

	t.Run(
		"BothAllDevices",
		func(t *testing.T) {
			t1 := lifxlan.AllDevices
			t2 := lifxlan.AllDevices
			if !t1.Matches(t2) {
				t.Errorf("%v should match %v", t1, t2)
			}
			if !t2.Matches(t1) {
				t.Errorf("%v should match %v", t2, t1)
			}
		},
	)
}

func TestParseTarget(t *testing.T) {
	t.Run(
		"EmptyAllDevices",
		func(t *testing.T) {
			s := ""
			target, err := lifxlan.ParseTarget(s)
			if err != nil {
				t.Fatal(err)
			}
			if target != lifxlan.AllDevices {
				t.Errorf(
					"ParseTarget(%q) expected %v, got %v",
					s,
					lifxlan.AllDevices,
					target,
				)
			}
		},
	)

	t.Run(
		"AllDevices",
		func(t *testing.T) {
			s := "00:00:00:00:00:00"
			target, err := lifxlan.ParseTarget(s)
			if err != nil {
				t.Fatal(err)
			}
			if target != lifxlan.AllDevices {
				t.Errorf(
					"ParseTarget(%q) expected %v, got %v",
					s,
					lifxlan.AllDevices,
					target,
				)
			}
		},
	)

	t.Run(
		"NormalDevice",
		func(t *testing.T) {
			s := "01:23:45:67:89:ab"
			target, err := lifxlan.ParseTarget(s)
			if err != nil {
				t.Fatal(err)
			}
			if target.String() != s {
				t.Errorf(
					"ParseTarget(%q) expected %v, got %v",
					s,
					lifxlan.AllDevices,
					target,
				)
			}
		},
	)
}

func ExampleTarget_Set() {
	var target lifxlan.Target
	flag.Var(
		&target,
		"target",
		"The MAC address of the target device. Empty value means any device.",
	)
	flag.Parse()
}
