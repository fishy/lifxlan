package relay_test

import (
	"context"
	"net"
	"testing"
	"time"

	"go.yhsif.com/lifxlan"
	"go.yhsif.com/lifxlan/mock"
	"go.yhsif.com/lifxlan/relay"
)

func TestWrap(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	const msg = relay.GetRPower
	const timeout = time.Millisecond * 200

	service, device := mock.StartService(t)

	t.Run(
		"Normal",
		func(t *testing.T) {
			service.RawStateRPowerPayload = &relay.RawStateRPowerPayload{
				Index: 0,
				Level: lifxlan.PowerOff,
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			if _, err := relay.Wrap(ctx, device, false); err != nil {
				t.Errorf("Expected successful wrapping, got: %v", err)
			}
		},
	)

	t.Run(
		"StateUnhandled",
		func(t *testing.T) {
			service.Handlers[msg] = mock.StateUnhandledHandler(msg)

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			if _, err := relay.Wrap(ctx, device, false); err == nil {
				t.Error("Expected Wrap to return error, got nil")
			} else {
				t.Logf("Got error: %v", err)
			}
		},
	)

	t.Run(
		"NoResponse",
		func(t *testing.T) {
			service.Handlers[msg] = func(*mock.Service, net.PacketConn, net.Addr, *lifxlan.Response) {}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			if _, err := relay.Wrap(ctx, device, false); err == nil {
				t.Error("Expected Wrap to return error, got nil")
			} else {
				t.Logf("Got error: %v", err)
			}
		},
	)
}
