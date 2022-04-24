package relay

import (
	"bytes"
	"context"
	"encoding/binary"
	"net"

	"go.yhsif.com/lifxlan"
)

// RawGetRPowerPayload defines the struct to be used for encoding and decoding.
//
// https://lan.developer.lifx.com/docs/querying-the-device-for-data#getrpower---packet-816
type RawGetRPowerPayload struct {
	Index uint8
}

// RawStateRPowerPayload defines the struct to be used for encoding and decoding.
//
// https://lan.developer.lifx.com/docs/information-messages#staterpower---packet-818
type RawStateRPowerPayload struct {
	Index uint8
	Level lifxlan.Power
}

func (rd *device) GetRPower(ctx context.Context, conn net.Conn, index uint8) (lifxlan.Power, error) {
	if ctx.Err() != nil {
		return 0, ctx.Err()
	}

	if conn == nil {
		newConn, err := rd.Dial()
		if err != nil {
			return 0, err
		}
		defer newConn.Close()
		conn = newConn

		if ctx.Err() != nil {
			return 0, ctx.Err()
		}
	}

	seq, err := rd.Send(
		ctx,
		conn,
		0, // flags
		GetRPower,
		&RawGetRPowerPayload{
			Index: index,
		},
	)
	if err != nil {
		return 0, err
	}

	for {
		resp, err := lifxlan.ReadNextResponse(ctx, conn)
		if err != nil {
			return 0, err
		}
		if resp.Sequence != seq || resp.Source != rd.Source() {
			continue
		}
		if resp.Message != StateRPower {
			continue
		}

		var raw RawStateRPowerPayload
		r := bytes.NewReader(resp.Payload)
		if err := binary.Read(r, binary.LittleEndian, &raw); err != nil {
			return 0, err
		}

		return raw.Level, nil
	}
}

// RawSetRPowerPayload defines the struct to be used for encoding and decoding.
//
// https://lan.developer.lifx.com/docs/changing-a-device#setrpower---packet-817
type RawSetRPowerPayload struct {
	Index uint8
	Level lifxlan.Power
}

func (rd *device) SetRPower(
	ctx context.Context,
	conn net.Conn,
	index uint8,
	power lifxlan.Power,
	ack bool,
) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if conn == nil {
		newConn, err := rd.Dial()
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

	seq, err := rd.Send(
		ctx,
		conn,
		flags,
		SetRPower,
		&RawSetRPowerPayload{
			Index: index,
			Level: power,
		},
	)
	if err != nil {
		return err
	}

	if ack {
		if err := lifxlan.WaitForAcks(ctx, conn, rd.Source(), seq); err != nil {
			return err
		}
	}

	return nil
}
