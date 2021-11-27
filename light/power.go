package light

import (
	"context"
	"net"
	"time"

	"go.yhsif.com/lifxlan"
)

// RawSetLightPowerPayload defines the struct to be used for encoding and decoding.
//
// https://lan.developer.lifx.com/docs/changing-a-device#setlightpower---packet-117
type RawSetLightPowerPayload struct {
	Level    lifxlan.Power
	Duration lifxlan.TransitionTime
}

func (ld *device) SetLightPower(
	ctx context.Context,
	conn net.Conn,
	power lifxlan.Power,
	transition time.Duration,
	ack bool,
) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if conn == nil {
		newConn, err := ld.Dial()
		if err != nil {
			return err
		}
		defer newConn.Close()
		conn = newConn

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	var flags lifxlan.AckResFlag
	if ack {
		flags |= lifxlan.FlagAckRequired
	}

	// Send
	seq, err := ld.Send(
		ctx,
		conn,
		flags,
		SetLightPower,
		&RawSetLightPowerPayload{
			Level:    power,
			Duration: lifxlan.ConvertDuration(transition),
		},
	)
	if err != nil {
		return err
	}

	if ack {
		return lifxlan.WaitForAcks(ctx, conn, ld.Source(), seq)
	}
	return nil
}
