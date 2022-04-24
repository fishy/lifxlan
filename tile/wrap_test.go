package tile_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"go.yhsif.com/lifxlan"
	"go.yhsif.com/lifxlan/light"
	"go.yhsif.com/lifxlan/mock"
	"go.yhsif.com/lifxlan/tile"
)

func mockProductMap(t *testing.T) {
	t.Helper()

	backupProductMap := lifxlan.ProductMap
	t.Cleanup(func() {
		lifxlan.ProductMap = backupProductMap
	})

	lifxlan.ProductMap = map[uint64]lifxlan.Product{
		lifxlan.ProductMapKey(1, 1): {
			ProductName: "Foo",
			Features: lifxlan.Features{
				Color: lifxlan.OptionalBoolPtr(true),
				Chain: lifxlan.OptionalBoolPtr(true),
			},
		},
		lifxlan.ProductMapKey(1, 2): {
			ProductName: "Boo",
			Features: lifxlan.Features{
				Color:            lifxlan.OptionalBoolPtr(true),
				Chain:            lifxlan.OptionalBoolPtr(true),
				TemperatureRange: lifxlan.TemperatureRange{100, 200},
			},
		},
	}
}

func TestWrap(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	mockProductMap(t)

	const timeout = time.Millisecond * 200

	var label lifxlan.Label
	label.Set("foo")

	service, device := mock.StartService(t)
	service.RawStatePayload = &light.RawStatePayload{
		Label: label,
	}

	rawEmpty := &tile.RawStateDeviceChainPayload{}

	expectedVersion1 := lifxlan.HardwareVersion{
		VendorID:        1,
		ProductID:       1,
		HardwareVersion: 1,
	}
	rawTile1 := tile.RawTileDevice{
		Width:           8,
		Height:          8,
		HardwareVersion: expectedVersion1,
	}
	rawOne1 := &tile.RawStateDeviceChainPayload{
		TotalCount: 1,
	}
	rawOne1.TileDevices[0] = rawTile1

	expectedVersion2 := lifxlan.HardwareVersion{
		VendorID:        1,
		ProductID:       2,
		HardwareVersion: 1,
	}
	rawTile2 := tile.RawTileDevice{
		Width:           8,
		Height:          8,
		HardwareVersion: expectedVersion2,
	}
	rawOne2 := &tile.RawStateDeviceChainPayload{
		TotalCount: 1,
	}
	rawOne2.TileDevices[0] = rawTile2

	t.Run(
		"EmptyChain",
		func(t *testing.T) {
			service.RawStateDeviceChainPayload = rawEmpty

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			if _, err := tile.Wrap(ctx, device, false); err == nil {
				t.Error("Expected error when no tiles in device, got nil")
			}
		},
	)

	t.Run(
		"OneTile",
		func(t *testing.T) {
			service.RawStateDeviceChainPayload = rawOne1

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			td, err := tile.Wrap(ctx, device, false)
			if err != nil {
				t.Fatal(err)
			}

			if td.Width() != 8 || td.Height() != 8 {
				t.Errorf("Got wrong size: %dx%d", td.Width(), td.Height())
			}

			if !reflect.DeepEqual(*td.HardwareVersion(), expectedVersion1) {
				t.Errorf(
					"td.HardwareVersion expected %v, got %v",
					expectedVersion1,
					td.HardwareVersion(),
				)
			}

			if !reflect.DeepEqual(*device.HardwareVersion(), expectedVersion1) {
				t.Errorf(
					"device.HardwareVersion expected %v, got %v",
					expectedVersion1,
					device.HardwareVersion(),
				)
			}

			service.RawStateDeviceChainPayload = rawOne2

			t.Run(
				"NoForce",
				func(t *testing.T) {
					ctx, cancel := context.WithTimeout(context.Background(), timeout)
					defer cancel()

					td, err := tile.Wrap(ctx, td, false)
					if err != nil {
						t.Fatal(err)
					}

					if !reflect.DeepEqual(*td.HardwareVersion(), expectedVersion1) {
						t.Errorf(
							"td.HardwareVersion expected %v, got %v",
							expectedVersion1,
							td.HardwareVersion(),
						)
					}

					if !reflect.DeepEqual(*device.HardwareVersion(), expectedVersion1) {
						t.Errorf(
							"device.HardwareVersion expected %v, got %v",
							expectedVersion1,
							device.HardwareVersion(),
						)
					}
				},
			)

			t.Run(
				"Force",
				func(t *testing.T) {
					ctx, cancel := context.WithTimeout(context.Background(), timeout)
					defer cancel()

					td, err := tile.Wrap(ctx, td, true)
					if err != nil {
						t.Fatal(err)
					}

					if !reflect.DeepEqual(*td.HardwareVersion(), expectedVersion2) {
						t.Errorf(
							"td.HardwareVersion expected %v, got %v",
							expectedVersion2,
							td.HardwareVersion(),
						)
					}

					if !reflect.DeepEqual(*device.HardwareVersion(), expectedVersion2) {
						t.Errorf(
							"device.HardwareVersion expected %v, got %v",
							expectedVersion2,
							device.HardwareVersion(),
						)
					}
				},
			)
		},
	)
}
