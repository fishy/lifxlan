package lifxlan

import (
	"bytes"
	"context"
	"encoding/binary"
	"net"
)

// Power is the raw power level value in messages.
//
// https://lan.developer.lifx.com/docs/device-messages#section-power-level
type Power uint16

// Power values.
const (
	PowerOn  Power = 65535
	PowerOff Power = 0
)

// On returns whether this power level value represents on state.
func (p Power) On() bool {
	return p != PowerOff
}

func (p Power) String() string {
	if p.On() {
		return "on"
	}
	return "off"
}

// RawStatePowerPayload defines the struct to be used for encoding and decoding.
//
// https://lan.developer.lifx.com/docs/device-messages#statepower---22
type RawStatePowerPayload struct {
	Level Power
}

func (d *device) GetPower(ctx context.Context, conn net.Conn) (Power, error) {
	select {
	default:
	case <-ctx.Done():
		return 0, ctx.Err()
	}

	if conn == nil {
		newConn, err := d.Dial()
		if err != nil {
			return 0, err
		}
		defer newConn.Close()
		conn = newConn

		select {
		default:
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}

	seq, err := d.Send(
		ctx,
		conn,
		0, // flags
		GetPower,
		nil, // payload
	)
	if err != nil {
		return 0, err
	}

	for {
		resp, err := ReadNextResponse(ctx, conn)
		if err != nil {
			return 0, err
		}
		if resp.Sequence != seq || resp.Source != d.Source() {
			continue
		}
		if resp.Message != StatePower {
			continue
		}

		var raw RawStatePowerPayload
		r := bytes.NewReader(resp.Payload)
		if err := binary.Read(r, binary.LittleEndian, &raw); err != nil {
			return 0, err
		}

		return raw.Level, nil
	}
}

// RawSetPowerPayload defines the struct to be used for encoding and decoding.
//
// https://lan.developer.lifx.com/docs/device-messages#setpower---21
type RawSetPowerPayload struct {
	Level Power
}

func (d *device) SetPower(
	ctx context.Context,
	conn net.Conn,
	power Power,
	ack bool,
) error {
	select {
	default:
	case <-ctx.Done():
		return ctx.Err()
	}

	if conn == nil {
		newConn, err := d.Dial()
		if err != nil {
			return err
		}
		defer newConn.Close()
		conn = newConn

		select {
		default:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	var flags AckResFlag
	if ack {
		flags |= FlagAckRequired
	}

	seq, err := d.Send(
		ctx,
		conn,
		flags,
		SetPower,
		&RawSetPowerPayload{
			Level: power,
		},
	)
	if err != nil {
		return err
	}

	if ack {
		if err := WaitForAcks(ctx, conn, d.Source(), seq); err != nil {
			return err
		}
	}

	return nil
}
