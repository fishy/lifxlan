package light_test

import (
	"context"
	"testing"
	"time"

	"go.yhsif.com/lifxlan"
	"go.yhsif.com/lifxlan/light"
	"go.yhsif.com/lifxlan/mock"
)

func TestWrap(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	const timeout = time.Millisecond * 200

	var label lifxlan.Label
	label.Set("foo")

	service, device := mock.StartService(t)

	t.Run(
		"Normal",
		func(t *testing.T) {
			service.RawStatePayload = &light.RawStatePayload{
				Label: label,
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			if _, err := light.Wrap(ctx, device, false); err != nil {
				t.Errorf("Expected successful wrapping, got: %v", err)
			}
		},
	)

	t.Run(
		"StateUnhandled",
		func(t *testing.T) {
			const msg = light.Get

			service.Handlers[msg] = mock.StateUnhandledHandler(msg)

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			if _, err := light.Wrap(ctx, device, false); err == nil {
				t.Error("Expected Wrap to return error, got nil")
			} else {
				t.Logf("Got error: %v", err)
			}
		},
	)
}
