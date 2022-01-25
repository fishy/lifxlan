package lifxlan_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"go.yhsif.com/lifxlan"
	"go.yhsif.com/lifxlan/mock"
)

func TestEcho(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	const timeout = time.Millisecond * 200

	seed := int64(0xdeadbeef)
	rand.Seed(seed)
	expectedSlice := make([]byte, 0, 64)
	rand.Read(expectedSlice)
	rand.Seed(seed)

	var expected [64]byte
	copy(expected[:], expectedSlice)

	service, device := mock.StartService(t)
	defer service.Stop()
	service.RawEchoResponsePayload = &lifxlan.RawEchoResponsePayload{
		Echoing: expected,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := device.Echo(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
}
