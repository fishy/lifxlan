package lifxlan_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/fishy/lifxlan"
	"github.com/fishy/lifxlan/mock"
)

func mockProductMap(t *testing.T) {
	t.Helper()
	lifxlan.ProductMap = map[uint64]lifxlan.ParsedHardwareVersion{
		lifxlan.ProductMapKey(1, 1): {
			ProductName: "Foo",
			Color:       true,
		},
	}
}

func TestVersion(t *testing.T) {
	mockProductMap(t)

	t.Run(
		"Found",
		func(t *testing.T) {
			raw := &lifxlan.RawHardwareVersion{
				VendorID:        1,
				ProductID:       1,
				HardwareVersion: 1,
			}
			expectedParsed := lifxlan.ParsedHardwareVersion{
				ProductName: "Foo",
				Color:       true,
				Raw:         *raw,
			}
			parsed := raw.Parse()
			if !reflect.DeepEqual(*parsed, expectedParsed) {
				t.Errorf("Parse expected %+v, got %+v", expectedParsed, parsed)
			}
			expectedStr := "Foo(1, 1, 1)"
			s := raw.String()
			if s != expectedStr {
				t.Errorf("String expected %q, got %q", expectedStr, s)
			}
		},
	)

	t.Run(
		"NotFound",
		func(t *testing.T) {
			raw := &lifxlan.RawHardwareVersion{
				VendorID:        1,
				ProductID:       2,
				HardwareVersion: 1,
			}
			parsed := raw.Parse()
			if parsed != nil {
				t.Errorf("Parse expected nil, got %+v", parsed)
			}
			expectedStr := "(1, 2, 1)"
			s := raw.String()
			if s != expectedStr {
				t.Errorf("String expected %q, got %q", expectedStr, s)
			}
		},
	)
}

func TestEmptyHardwareVersion(t *testing.T) {
	var version lifxlan.RawHardwareVersion
	s := version.String()
	if s != lifxlan.EmptyHardwareVersion {
		t.Errorf("Expected %q, got %q", lifxlan.EmptyHardwareVersion, s)
	}
}

func TestGetHardwareVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	const timeout = time.Millisecond * 200

	mockProductMap(t)
	expected := lifxlan.RawHardwareVersion{
		VendorID:        1,
		ProductID:       1,
		HardwareVersion: 1,
	}

	service, device := mock.StartService(t)
	defer service.Stop()
	service.RawStateVersionPayload = &lifxlan.RawStateVersionPayload{
		Version: expected,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := device.GetHardwareVersion(ctx, nil); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*device.HardwareVersion(), expected) {
		t.Errorf(
			"HardwareVersion expected %v, got %v",
			expected,
			device.HardwareVersion(),
		)
	}
}
