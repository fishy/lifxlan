package relay_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"net"
	"testing"
	"time"

	"go.yhsif.com/lifxlan"
	"go.yhsif.com/lifxlan/mock"
	"go.yhsif.com/lifxlan/relay"
)

func TestGetRPower(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	const timeout = time.Millisecond * 200

	const (
		expected = lifxlan.PowerOn
		index    = 0
	)

	service, device := mock.StartService(t)
	service.RawStateRPowerPayload = &relay.RawStateRPowerPayload{
		Index: index,
		Level: expected,
	}

	rd, err := func() (relay.Device, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return relay.Wrap(ctx, device, false)
	}()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	power, err := rd.GetRPower(ctx, nil, index)
	if err != nil {
		t.Fatal(err)
	}
	if expected != power {
		t.Errorf(
			"Power expected %v(%d), got %v(%d)",
			expected,
			expected,
			power,
			power,
		)
	}
}

func TestSetRPower(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	const timeout = time.Millisecond * 200

	const (
		index    = 1
		expected = lifxlan.PowerOff
	)

	service, device := mock.StartService(t)
	service.RawStateRPowerPayload = &relay.RawStateRPowerPayload{
		Index: index,
	}

	rd, err := func() (relay.Device, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return relay.Wrap(ctx, device, false)
	}()
	if err != nil {
		t.Fatal(err)
	}

	var called bool

	service.Handlers[relay.SetRPower] = func(
		_ *mock.Service,
		_ net.PacketConn,
		_ net.Addr,
		orig *lifxlan.Response,
	) {
		called = true
		var raw relay.RawSetRPowerPayload
		r := bytes.NewReader(orig.Payload)
		if err := binary.Read(r, binary.LittleEndian, &raw); err != nil {
			t.Fatal(err)
		}
		if expected != raw.Level {
			t.Errorf(
				"Power expected %v(%d), got %v(%d)",
				expected,
				expected,
				raw.Level,
				raw.Level,
			)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := rd.SetRPower(ctx, nil, index, expected, true); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("SetRPower message not received.")
	}
}
