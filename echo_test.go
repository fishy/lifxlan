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

	err := device.Echo(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
}
