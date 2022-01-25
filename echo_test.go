package lifxlan_test

import (
	"context"
	"testing"
	"time"

	"go.yhsif.com/lifxlan/mock"
)

func TestEcho(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	const timeout = time.Millisecond * 200

	service, device := mock.StartService(t)
	defer service.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	t.Run(
		"UserPayload",
		func(t *testing.T) {
			err := device.Echo(ctx, nil, []byte("payload"))
			if err != nil {
				t.Fatal(err)
			}
		},
	)

	t.Run(
		"OversizePayload",
		func(t *testing.T) {
			err := device.Echo(ctx, nil, []byte("this is a big string that's longer than the allowed 64 bytes for the echo payload"))
			if err != nil {
				t.Fatal(err)
			}
		},
	)

	t.Run(
		"DefaultPayload",
		func(t *testing.T) {
			err := device.Echo(ctx, nil, nil)
			if err != nil {
				t.Fatal(err)
			}
		},
	)
}
