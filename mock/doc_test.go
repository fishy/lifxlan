package mock_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"net"
	"testing"
	"time"

	"github.com/fishy/lifxlan"
	"github.com/fishy/lifxlan/mock"
	"github.com/fishy/lifxlan/tile"
)

// This example demonstrates how to mock a response in test code.
func Example_testGetLabel() {
	var t *testing.T
	t.Run(
		"GetLabel",
		func(t *testing.T) {
			if testing.Short() {
				t.Skip("skipping test in short mode.")
			}

			const timeout = time.Millisecond * 200

			var expected lifxlan.Label
			expected.Set("foo")

			service, device := mock.StartService(t)
			defer service.Stop()
			// This is the payload to be returned by the mock service.
			service.RawStateLabelPayload = &lifxlan.RawStateLabelPayload{
				Label: expected,
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			if err := device.GetLabel(ctx, nil); err != nil {
				t.Fatal(err)
			}
			if device.Label().String() != expected.String() {
				t.Errorf("Label expected %v, got %v", expected, device.Label())
			}
		},
	)
}

// This example demonstrates how to mock a response with custom HandlerFunc.
func Example_testGetLabelWithHandlerFunc() {
	var t *testing.T
	t.Run(
		"GetLabel",
		func(t *testing.T) {
			if testing.Short() {
				t.Skip("skipping test in short mode.")
			}

			const timeout = time.Millisecond * 200

			var expected lifxlan.Label
			expected.Set("foo")

			service, device := mock.StartService(t)
			defer service.Stop()

			// This defines the handler for GetLabel messages.
			service.Handlers[lifxlan.GetLabel] = func(
				s *mock.Service,
				conn net.PacketConn,
				addr net.Addr,
				orig *lifxlan.Response,
			) {
				buf := new(bytes.Buffer)
				if err := binary.Write(
					buf,
					binary.LittleEndian,
					expected,
				); err != nil {
					s.TB.Log(err)
					return
				}
				s.Reply(conn, addr, orig, lifxlan.StateLabel, buf.Bytes())
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			if err := device.GetLabel(ctx, nil); err != nil {
				t.Fatal(err)
			}
			if device.Label().String() != expected.String() {
				t.Errorf("Label expected %v, got %v", expected, device.Label())
			}
		},
	)
}

// This example demonstrates how to mock a not enough acks situation in test
// code.
func Example_testNotEnoughAcks() {
	var t *testing.T
	t.Run(
		"SetColors",
		func(t *testing.T) {
			if testing.Short() {
				t.Skip("skipping test in short mode.")
			}

			const timeout = time.Millisecond * 200

			service, device := mock.StartService(t)
			defer service.Stop()

			rawTile1 := tile.RawTileDevice{
				Width:  8,
				Height: 8,
			}
			rawTile2 := tile.RawTileDevice{
				UserX:  1,
				Width:  8,
				Height: 8,
			}
			rawChain := &tile.RawStateDeviceChainPayload{
				TotalCount: 2,
			}
			rawChain.TileDevices[0] = rawTile1
			rawChain.TileDevices[1] = rawTile2
			service.RawStateDeviceChainPayload = rawChain

			td, err := func() (tile.Device, error) {
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()
				return tile.Wrap(ctx, device, false)
			}()
			if err != nil {
				t.Fatal(err)
			}
			if td == nil {
				t.Fatal("Can't mock tile device.")
			}

			t.Run(
				"NotEnoughAcks",
				func(t *testing.T) {
					// The SetColors function will expect 2 acks.
					service.AcksToDrop = 1

					ctx, cancel := context.WithTimeout(context.Background(), timeout)
					defer cancel()

					if err := td.SetColors(ctx, nil, nil, 0, true); err == nil {
						t.Error("Expected error when not enough acks returned, got nil")
					}
				},
			)
		},
	)
}

// Eliminates the "no tests to run" warning.
func TestEmpty(t *testing.T) {
}
