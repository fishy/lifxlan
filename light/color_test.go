package light_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"net"
	"reflect"
	"testing"
	"time"

	"go.yhsif.com/lifxlan"
	"go.yhsif.com/lifxlan/light"
	"go.yhsif.com/lifxlan/mock"
)

func mockProductMap(t *testing.T) {
	t.Helper()
	lifxlan.ProductMap = map[uint64]lifxlan.ParsedHardwareVersion{
		lifxlan.ProductMapKey(1, 1): {
			ProductName: "Bar",
			Color:       true,
			Chain:       true,
			MinKelvin:   100,
			MaxKelvin:   200,
		},
	}
}

func TestSetColor(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	mockProductMap(t)

	version := lifxlan.HardwareVersion{
		VendorID:        1,
		ProductID:       1,
		HardwareVersion: 1,
	}

	const timeout = time.Millisecond * 200

	var label lifxlan.Label
	label.Set("foo")

	service, device := mock.StartService(t)
	defer service.Stop()
	service.RawStatePayload = &light.RawStatePayload{
		Label: label,
	}

	ld, err := func() (light.Device, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return light.Wrap(ctx, device, false)
	}()
	if err != nil {
		t.Fatal(err)
	}

	color := lifxlan.Color{
		Hue:        1,
		Saturation: 2,
		Brightness: 3,
		Kelvin:     0,
	}

	t.Run(
		"GetColor",
		func(t *testing.T) {
			service.RawStatePayload = &light.RawStatePayload{
				Color: color,
				Label: label,
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			ret, err := ld.GetColor(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(*ret, color) {
				t.Errorf("Expected color %+v, got %+v", color, *ret)
			}
			if gotLabel := ld.Label().String(); gotLabel != label.String() {
				t.Errorf("Expected label %q, got %q", label.String(), gotLabel)
			}
		},
	)

	t.Run(
		"SetColor",
		func(t *testing.T) {
			t.Run(
				"NoAck",
				func(t *testing.T) {
					service.AcksToDrop = 1

					ctx, cancel := context.WithTimeout(context.Background(), timeout)
					defer cancel()

					if err := ld.SetColor(ctx, nil, &color, 0, true); err == nil {
						t.Error("Expected error when not getting ack, got nil")
					}
				},
			)

			t.Run(
				"Normal",
				func(t *testing.T) {
					service.AcksToDrop = 0

					*ld.HardwareVersion() = version

					service.Handlers[light.SetColor] = func(
						_ *mock.Service,
						_ net.PacketConn,
						_ net.Addr,
						orig *lifxlan.Response,
					) {
						var raw light.RawSetColorPayload
						r := bytes.NewReader(orig.Payload)
						if err := binary.Read(r, binary.LittleEndian, &raw); err != nil {
							t.Fatal(err)
						}
						parsed := ld.HardwareVersion().Parse()
						if parsed == nil {
							t.Fatal("No hardware version cached")
						}
						k := raw.Color.Kelvin
						if k < parsed.MinKelvin || k > parsed.MaxKelvin {
							t.Errorf("Color not sanitized: %+v", raw.Color)
						}
					}

					ctx, cancel := context.WithTimeout(context.Background(), timeout)
					defer cancel()

					if err := ld.SetColor(ctx, nil, &color, 0, true); err != nil {
						t.Error(err)
					}
				},
			)
		},
	)
}
