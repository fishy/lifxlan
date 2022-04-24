package lifxlan_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"math"
	"math/rand"
	"net"
	"testing"
	"testing/quick"
	"time"

	"go.yhsif.com/lifxlan"
	"go.yhsif.com/lifxlan/mock"
)

func TestPower(t *testing.T) {
	t.Run(
		"Off",
		func(t *testing.T) {
			expected := "off"
			power := lifxlan.PowerOff
			if power.On() {
				t.Errorf("%v.On() should return false.", power)
			}
			s := power.String()
			if s != expected {
				t.Errorf("Power(%d).String() expected %q, got %q", power, expected, s)
			}
		},
	)

	t.Run(
		"On",
		func(t *testing.T) {
			expected := "on"
			power := lifxlan.PowerOn
			if !power.On() {
				t.Errorf("Power(%d).On() should return true.", power)
			}
			s := power.String()
			if s != expected {
				t.Errorf("Power(%d).String() expected %q, got %q", power, expected, s)
			}
		},
	)

	t.Run(
		"RandomOn",
		func(t *testing.T) {
			// Seed rander
			now := time.Now()
			rander := rand.New(rand.NewSource(now.Unix() + int64(now.Nanosecond())))

			expected := "on"
			var n int
			f := func() bool {
				n++
				power := lifxlan.Power(rander.Intn(math.MaxUint16-1) + 1)
				pass := true
				if !power.On() {
					pass = false
					t.Logf("Power(%d).On() should return true.", power)
				}
				s := power.String()
				if s != expected {
					pass = false
					t.Logf("Power(%d).String() expected %q, got %q", power, expected, s)
				}
				return pass
			}
			if err := quick.Check(f, nil); err != nil {
				t.Error(err)
			}
			t.Logf("quick did %d checks", n)
		},
	)
}

func TestGetPower(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	const timeout = time.Millisecond * 200

	expected := lifxlan.PowerOn

	service, device := mock.StartService(t)
	service.RawStatePowerPayload = &lifxlan.RawStatePowerPayload{
		Level: expected,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	power, err := device.GetPower(ctx, nil)
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

func TestSetPower(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	const timeout = time.Millisecond * 200

	expected := lifxlan.PowerOff

	service, device := mock.StartService(t)

	var called bool

	service.Handlers[lifxlan.SetPower] = func(
		_ *mock.Service,
		_ net.PacketConn,
		_ net.Addr,
		orig *lifxlan.Response,
	) {
		called = true
		var raw lifxlan.RawSetPowerPayload
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

	if err := device.SetPower(ctx, nil, expected, true); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("SetPower message not received.")
	}
}
