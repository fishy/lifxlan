package lifxlan_test

import (
	"fmt"
	"testing"

	"github.com/fishy/lifxlan"
)

func TestTargetString(t *testing.T) {
	target := lifxlan.AllDevices
	str := fmt.Sprintf("%v", target)
	expected := "00:00:00:00:00:00"
	if str != expected {
		t.Errorf("Expected %s, got %s", expected, str)
	}

	target = 0xffffffd573d0
	str = fmt.Sprintf("%v", target)
	expected = "d0:73:d5:ff:ff:ff"
	if str != expected {
		t.Errorf("Expected %s, got %s", expected, str)
	}
}
